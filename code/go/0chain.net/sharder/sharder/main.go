package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	_ "net/http/pprof"

	"0chain.net/chaincore/block"
	"0chain.net/chaincore/chain"
	"0chain.net/chaincore/client"
	"0chain.net/chaincore/config"
	"0chain.net/chaincore/diagnostics"
	"0chain.net/chaincore/node"
	"0chain.net/chaincore/round"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/build"
	"0chain.net/core/common"
	"0chain.net/core/ememorystore"
	"0chain.net/core/encryption"
	"0chain.net/core/logging"
	. "0chain.net/core/logging"
	"0chain.net/core/memorystore"
	"0chain.net/core/persistencestore"
	"0chain.net/sharder"
	"0chain.net/sharder/blockstore"
	"0chain.net/smartcontract/setupsc"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	deploymentMode := flag.Int("deployment_mode", 2, "deployment_mode")
	keysFile := flag.String("keys_file", "", "keys_file")
	nodesFile := flag.String("nodes_file", "", "nodes_file (deprecated)")
	maxDelay := flag.Int("max_delay", 0, "max_delay (deprecated)")
	flag.Parse()
	config.Configuration.DeploymentMode = byte(*deploymentMode)
	config.SetupDefaultConfig()
	config.SetupConfig()
	config.SetupSmartContractConfig()

	if config.Development() {
		logging.InitLogging("development")
	} else {
		logging.InitLogging("production")
	}

	config.Configuration.ChainID = viper.GetString("server_chain.id")
	config.Configuration.MaxDelay = *maxDelay

	reader, err := os.Open(*keysFile)
	if err != nil {
		panic(err)
	}

	config.SetServerChainID(config.Configuration.ChainID)
	common.SetupRootContext(node.GetNodeContext())
	ctx := common.GetRootContext()
	initEntities()
	serverChain := chain.NewChainFromConfig()
	signatureScheme := serverChain.GetSignatureScheme()
	err = signatureScheme.ReadKeys(reader)
	if err != nil {
		Logger.Panic("Error reading keys file")
	}
	node.Self.SetSignatureScheme(signatureScheme)
	reader.Close()

	sharder.SetupSharderChain(serverChain)
	sc := sharder.GetSharderChain()
	chain.SetServerChain(serverChain)
	chain.SetNetworkRelayTime(viper.GetDuration("network.relay_time") * time.Millisecond)
	node.ReadConfig()

	nodesConfigFile := viper.GetString("network.nodes_file")
	if nodesConfigFile == "" {
		nodesConfigFile = *nodesFile
	}
	if nodesConfigFile == "" {
		panic("Please specify --nodes_file file.txt option with a file.txt containing nodes including self")
	}
	if strings.HasSuffix(nodesConfigFile, "txt") {
		reader, err = os.Open(nodesConfigFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
		node.ReadNodes(reader, serverChain.Miners, serverChain.Sharders, serverChain.Blobbers)
		reader.Close()
	} else {
		sc.ReadNodePools(nodesConfigFile)
		Logger.Info("nodes", zap.Int("miners", sc.Miners.Size()), zap.Int("sharders", sc.Sharders.Size()))
	}

	if node.Self.ID == "" {
		Logger.Panic("node definition for self node doesn't exist")
	}
	if node.Self.Type != node.NodeTypeSharder {
		Logger.Panic("node not configured as sharder")
	}

	if state.Debug() {
		chain.SetupStateLogger("/tmp/state.txt")
	}
	sc.SetupGenesisBlock(viper.GetString("server_chain.genesis_block.id"))

	mode := "main net"
	if config.Development() {
		mode = "development"
	} else if config.TestNet() {
		mode = "test net"
	}

	address := fmt.Sprintf(":%v", node.Self.Port)

	Logger.Info("Starting sharder", zap.String("git", build.GitCommit), zap.String("go_version", runtime.Version()), zap.Int("available_cpus", runtime.NumCPU()), zap.String("port", address))
	Logger.Info("Chain info", zap.String("chain_id", config.GetServerChainID()), zap.String("mode", mode))
	Logger.Info("Self identity", zap.Any("set_index", node.Self.Node.SetIndex), zap.Any("id", node.Self.Node.GetKey()))

	var server *http.Server
	if config.Development() {
		// No WriteTimeout setup to enable pprof
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	} else {
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}
	common.HandleShutdown(server)
	setupBlockStorageProvider()

	initWorkers(ctx)
	common.ConfigRateLimits()
	initN2NHandlers()
	initServer()
	initHandlers()

	Logger.Info("In Sharder rejoin")
	r, err := sc.GetMostRecentRoundFromDB(ctx)

	if err == nil {
		sc.CurrentRound = r.Number
		sc.AddRound(r)
		Logger.Info("bc-27 most recent round info", zap.Int64("roundNumber", r.Number), zap.String("blockHash", r.BlockHash))
		} else {
		Logger.Error("bc-27 error reading round data from db", zap.Error(err))
	}

	go syncUpRounds(ctx, r)
	
	Logger.Info("Ready to listen to the requests")
	chain.StartTime = time.Now().UTC()
	log.Fatal(server.ListenAndServe())
}

func initServer() {
	// TODO; when a new server is brought up, it needs to first download all the state before it can start accepting requests
}

func initHandlers() {
	if config.Development() {
		http.HandleFunc("/_hash", encryption.HashHandler)
		http.HandleFunc("/_sign", common.ToJSONResponse(encryption.SignHandler))
	}
	config.SetupHandlers()
	node.SetupHandlers()
	chain.SetupHandlers()
	block.SetupHandlers()
	sharder.SetupHandlers()
	diagnostics.SetupHandlers()
	chain.SetupStateHandlers()

	serverChain := chain.GetServerChain()
	serverChain.SetupNodeHandlers()
}

func initEntities() {
	memoryStorage := memorystore.GetStorageProvider()
	chain.SetupEntity(memoryStorage)
	block.SetupEntity(memoryStorage)

	round.SetupRoundSummaryDB()
	block.SetupBlockSummaryDB()
	ememoryStorage := ememorystore.GetStorageProvider()
	block.SetupBlockSummaryEntity(ememoryStorage)
	block.SetupStateChange(memoryStorage)
	state.SetupPartialState(memoryStorage)
	state.SetupStateNodes(memoryStorage)
	round.SetupEntity(ememoryStorage)

	client.SetupEntity(memoryStorage)
	transaction.SetupEntity(memoryStorage)

	persistencestore.InitSession()
	persistenceStorage := persistencestore.GetStorageProvider()
	transaction.SetupTxnSummaryEntity(persistenceStorage)
	transaction.SetupTxnConfirmationEntity(persistenceStorage)

	if config.DevConfiguration.SmartContract {
		setupsc.SetupSmartContracts()
	}
}

func initN2NHandlers() {
	node.SetupN2NHandlers()
	sharder.SetupM2SReceivers()
	sharder.SetupM2SResponders()
	chain.SetupX2MRequestors()
	sharder.SetupS2SRequestors()
	sharder.SetupS2SResponders()
}

func initWorkers(ctx context.Context) {
	serverChain := chain.GetServerChain()
	serverChain.SetupWorkers(ctx)
	sharder.SetupWorkers(ctx)
}

func syncUpRounds(ctx context.Context, r *round.Round) {
	sc := sharder.GetSharderChain()
	sc.Sharders.OneTimeStatusMonitor(ctx)
	Logger.Info("bc-27 get latest round from other sharders", zap.Int64("sharder_curr_round", r.Number))
	lr := sc.GetLatestRoundFromSharders(ctx, r.Number)
	if lr != nil && lr.Number > r.Number + 1 {
		sc.SetStatus(sharder.SharderSyncing)
		Logger.Info("#rejoin sharder status set to syncing")
		Logger.Info("bc-27 latest round from other sharder", zap.Int64("curr_round", r.Number), zap.Int64("latest_round_from_sharders", lr.Number))		
		ts := time.Now()
		sc.GetMissingRounds(ctx, lr.Number, r.Number)
		duration := time.Since(ts)
		Logger.Info("bc-27 duration for catching up with all missing rounds", zap.Duration("duration", duration))
		sc.SetStatus(sharder.SharderNormal)
		Logger.Info("#rejoin sharder status set to Normal")
	}
	go sc.BlockWorker(ctx)
}

func setupBlockStorageProvider() {
	blockStorageProvider := viper.GetString("server_chain.block.storage.provider")
	if blockStorageProvider == "" || blockStorageProvider == "blockstore.FSBlockStore" {
		blockstore.SetupStore(blockstore.NewFSBlockStore("data/blocks"))
	} else if blockStorageProvider == "blockstore.BlockDBStore" {
		blockstore.SetupStore(blockstore.NewBlockDBStore("data/blocksdb"))
	} else if blockStorageProvider == "blockstore.MultiBlockstore" {
		var bs = []blockstore.BlockStore{blockstore.NewFSBlockStore("data/blocks"), blockstore.NewBlockDBStore("data/blocksdb")}
		blockstore.SetupStore(blockstore.NewMultiBlockStore(bs))
	} else {
		panic(fmt.Sprintf("uknown block store provider - %v", blockStorageProvider))
	}
}
