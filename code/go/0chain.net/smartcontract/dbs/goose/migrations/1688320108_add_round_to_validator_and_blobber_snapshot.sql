-- +goose Up
-- +goose StatementBegin
ALTER TABLE validator_snapshots ADD COLUMN round BIGINT NOT NULL DEFAULT 0;
ALTER TABLE blobber_snapshots ADD COLUMN round BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE validator_snapshots DROP COLUMN round;
ALTER TABLE blobber_snapshots DROP COLUMN round;
-- +goose StatementEnd
