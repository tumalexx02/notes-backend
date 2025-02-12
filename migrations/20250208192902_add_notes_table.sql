-- +goose Up
-- +goose StatementBegin
CREATE TABLE notes (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  archived_at TIMESTAMP
);
CREATE INDEX idx_notes_updated_at ON notes (updated_at);
CREATE TABLE note_nodes (
  id SERIAL PRIMARY KEY,
  note_id INTEGER NOT NULL,
  "order" INTEGER NOT NULL,
  content_type TEXT NOT NULL,
  content TEXT NOT NULL,
  FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_notes_updated_at;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS note_nodes;
-- +goose StatementEnd
