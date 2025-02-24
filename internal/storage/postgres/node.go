package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models/note"
	"main/internal/storage"
)

func (s *Storage) AddNoteNode(noteId int, contentType string, content string) (int, error) {
	const op = "storage.postgres.CreateNoteNode"

	// begin transaction
	tx := s.db.MustBegin()

	// creating note node
	var id int

	err := tx.Get(&id, createBlankNoteNodeQuery, noteId, contentType, content)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// set updated_at field on note
	_, err = tx.Exec(setUpdatedAtQuery, noteId)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteNoteNode(id int) error {
	const op = "storage.postgres.DeleteNoteNode"

	// begin transaction
	tx := s.db.MustBegin()

	// deleting node with returning (note_id, order)
	var tempNoteNode note.NoteNode

	err := tx.Get(&tempNoteNode, deleteNoteNodeQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// getting note_id and order from deleted node
	var noteId, order = tempNoteNode.NoteId, tempNoteNode.Order

	// update all note nodes' order after deleted node
	_, err = tx.Exec(updateOrderAfterDeleteQuery, noteId, order)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// set updated_at field on note
	_, err = tx.Exec(setUpdatedAtQuery, noteId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateNoteNodeContent(id int, content string) error {
	const op = "storage.postgres.UpdateNoteNodeContent"

	// begin transaction
	tx := s.db.MustBegin()

	// updating note node with returning note_id
	var noteId int

	err := tx.Get(&noteId, updateNoteNodeContentQuery, id, content)
	if errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return storage.ErrNoteNodeNotFound
	}
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// set updated_at field on note
	_, err = tx.Exec(setUpdatedAtQuery, noteId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) IsUserNoteNodeOwner(userId string, noteNodeId int) (bool, error) {
	const op = "storage.postgres.IsUserNoteNodeOwner"

	var exists int

	err := s.db.Get(&exists, isUserNodeOwnerQuery, userId, noteNodeId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists == 1, nil
}

func (s *Storage) GetNodeById(id int) (note.NoteNode, error) {
	const op = "storage.postgres.GetNoteIdByNoteNodeId"

	// getting note id by note node id
	var node note.NoteNode

	err := s.db.Get(&node, getNoteIdByNoteNodeIdQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return node, storage.ErrNoteNodeNotFound
	}
	if err != nil {
		return node, fmt.Errorf("%s: %w", op, err)
	}

	return node, nil
}

func (s *Storage) GetAllNotesNodes(noteId int) ([]note.NoteNode, error) {
	const op = "storage.postgres.GetAllNotesNodes"

	var nodes []note.NoteNode

	err := s.db.Select(&nodes, getAllNotesNodesQuery, noteId)
	if errors.Is(err, sql.ErrNoRows) {
		return nodes, storage.ErrNoteNotFound
	}
	if err != nil {
		return nodes, fmt.Errorf("%s: %w", op, err)
	}

	return nodes, nil
}
