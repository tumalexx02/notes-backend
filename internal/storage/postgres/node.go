package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/storage"
)

const (
	createBlankNoteNodeQuery = `
		INSERT INTO note_nodes (note_id, "order", content_type, content) 
		VALUES ($1, (SELECT COUNT(*) FROM note_nodes WHERE note_id = $1), $2, $3)
		RETURNING id;
	`
	deleteNoteNodeQuery = `
		DELETE FROM note_nodes
		WHERE id = $1;
	`
	getOrderQuery = `
		SELECT "order"
		FROM note_nodes
		WHERE id = $1;
	`
	getNoteIdQuery = `
		SELECT note_id
		FROM note_nodes
		WHERE id = $1;
	`
	updateNoteNodesOrderQuery = `
		UPDATE note_nodes
		SET "order" = "order" - 1
		WHERE note_id = $1 AND "order" > $2;
	`
	updateNoteNodeContentQuery = `
		UPDATE note_nodes
		SET content = $2
		WHERE id = $1;
	`
)

func (s *Storage) AddNoteNode(noteId int, contentType string, content string) (int, error) {
	const op = "storage.postgres.CreateNoteNode"

	var id int

	err := s.db.Get(&id, createBlankNoteNodeQuery, noteId, contentType, content)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteNoteNode(id int) error {
	const op = "storage.postgres.DeleteNoteNode"

	tx := s.db.MustBegin()

	var order int
	err := tx.Get(&order, getOrderQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = tx.Rollback()
			return storage.ErrNoteNodeNotFound
		}
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	var noteId int
	err = tx.Get(&noteId, getNoteIdQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(deleteNoteNodeQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(updateNoteNodesOrderQuery, noteId, order)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateNoteNodeContent(id int, content string) error {
	const op = "storage.postgres.UpdateNoteNodeContent"

	_, err := s.db.Exec(updateNoteNodeContentQuery, id, content)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
