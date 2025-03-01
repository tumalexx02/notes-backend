-- +goose Up
-- +goose StatementBegin
ALTER TABLE notes 
ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT FALSE;
CREATE INDEX idx_notes_is_public ON notes (is_public);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_notes_is_public;
ALTER TABLE notes 
DROP COLUMN IF EXISTS is_public;
-- +goose StatementEnd
