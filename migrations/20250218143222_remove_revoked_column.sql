-- +goose Up
-- +goose StatementBegin
ALTER TABLE refresh_tokens
  DROP COLUMN IF EXISTS revoked;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE refresh_tokens
  ADD COLUMN revoked BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd
