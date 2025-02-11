package postgres

import (
	"fmt"
	"main/internal/storage"
)

const (
	createBlankNoteNodeQuery = `
		INSERT INTO note_nodes (note_id, "order", content_type, content) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	getNodesCountByIdQuery = `
		SELECT COUNT(*) FROM note_nodes
		WHERE note_id = $1;
	`
	deleteNoteNodeQuery = `
		DELETE FROM note_nodes
		WHERE id = $1;
	`
)

func (s *Storage) AddNoteNode(noteId int, contentType string, content string) (int, error) {
	const op = "storage.postgres.CreateNoteNode"

	var nodesCount int
	err := s.db.Get(&nodesCount, getNodesCountByIdQuery, noteId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int

	err = s.db.Get(&id, createBlankNoteNodeQuery, noteId, nodesCount, contentType, content)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteNoteNode(id int) error {
	const op = "storage.postgres.DeleteNoteNode"

	res, err := s.db.Exec(deleteNoteNodeQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrNoteNodeNotFound
	}

	return nil
}
