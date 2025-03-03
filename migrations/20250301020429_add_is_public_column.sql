-- +goose Up
-- +goose StatementBegin
ALTER TABLE notes 
ADD COLUMN public_id UUID DEFAULT NULL;
CREATE INDEX idx_notes_public_note ON notes (public_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_notes_public_note;
ALTER TABLE notes 
DROP COLUMN IF EXISTS public_id;
-- +goose StatementEnd
