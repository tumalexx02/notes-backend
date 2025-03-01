-- Active: 1739048967498@@127.0.0.1@5432@postgres
-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_note_nodes_note_id ON note_nodes (note_id);
CREATE INDEX idx_note_nodes_order ON note_nodes ("order");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_note_nodes_note_id;
DROP INDEX IF EXISTS idx_note_nodes_order;
-- +goose StatementEnd
