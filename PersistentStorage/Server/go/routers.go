/*
 * 0chain Services
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"ClientsDelete",
		strings.ToUpper("Delete"),
		"/clients",
		ClientsDelete,
	},

	Route{
		"ClientsGet",
		strings.ToUpper("Get"),
		"/clients",
		ClientsGet,
	},

	Route{
		"ClientsPatch",
		strings.ToUpper("Patch"),
		"/clients",
		ClientsPatch,
	},

	Route{
		"ClientsPost",
		strings.ToUpper("Post"),
		"/clients",
		ClientsPost,
	},

	Route{
		"TransactionsDelete",
		strings.ToUpper("Delete"),
		"/transactions",
		TransactionsDelete,
	},

	Route{
		"TransactionsGet",
		strings.ToUpper("Get"),
		"/transactions",
		TransactionsGet,
	},

	Route{
		"TransactionsPatch",
		strings.ToUpper("Patch"),
		"/transactions",
		TransactionsPatch,
	},

	Route{
		"TransactionsPost",
		strings.ToUpper("Post"),
		"/transactions",
		TransactionsPost,
	},
}
