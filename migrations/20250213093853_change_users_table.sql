-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
  ALTER COLUMN created_at SET DATA TYPE TIMESTAMP,
  ALTER COLUMN updated_at SET DATA TYPE TIMESTAMP;

ALTER TABLE users
  ADD COLUMN name TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
  DROP COLUMN IF EXISTS name;

ALTER TABLE users
  ALTER COLUMN created_at SET DATA TYPE TIMESTAMPTZ,
  ALTER COLUMN updated_at SET DATA TYPE TIMESTAMPTZ;
-- +goose StatementEnd