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
		WHERE id = $1
		RETURNING note_id;
	`
)

func (s *Storage) AddNoteNode(noteId int, contentType string, content string) (int, error) {
	const op = "storage.postgres.CreateNoteNode"

	tx := s.db.MustBegin()

	var id int

	err := tx.Get(&id, createBlankNoteNodeQuery, noteId, contentType, content)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(setUpdatedAtQuery, noteId)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
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

	_, err = tx.Exec(setUpdatedAtQuery, noteId)
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

	tx := s.db.MustBegin()

	var noteId int

	err := tx.Get(&noteId, updateNoteNodeContentQuery, id, content)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(setUpdatedAtQuery, noteId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// TODO: update order method
// old_order -> new_order
