package postgres

import (
	"fmt"
	"main/internal/models/note"
)

const (
	createNoteQuery = `
		INSERT INTO notes (title, user_id) 
		VALUES ($1, $2)
		RETURNING id;
	`
	createBlankNoteNodeQuery = `
		INSERT INTO note_nodes (note_id, "order", content_type, content) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	getNodesCountByIdQuery = `
		SELECT COUNT(*) FROM note_nodes
		WHERE note_id = $1;
	`
	getNoteQuery = `
		SELECT * FROM notes
		WHERE id = $1;
	`
	getNotesByUserIdQuery = `
		SELECT * FROM notes
		WHERE user_id = $1;
	`
	updateNoteTitleQuery = `
		UPDATE notes
		SET title = $2, updated_at = NOW()
		WHERE id = $1;
	`
	updateNoteNodeQuery = `
		UPDATE note_nodes
		SET content_type = $2, content = $3, updated_at = NOW()
		WHERE id = $1;
	`
	archiveNoteQuery = `
		UPDATE notes
		SET archived_at = NOW()
		WHERE id = $1;
	`
	deleteNoteQuery = `
		DELETE FROM notes
		WHERE id = $1;
	`
)

func (s *Storage) CreateNote(noteTitle string, userId string) (uint, error) {
	const op = "storage.postgres.CreateNote"

	var id uint

	err := s.db.QueryRow(createNoteQuery, noteTitle, userId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	blankTextNode := note.NoteNode{
		NoteId:      id,
		Order:       0,
		ContentType: note.ContentTypeText,
		Content:     "",
	}

	_, err = s.AddNoteNode(blankTextNode)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) AddNoteNode(newNoteNode note.NoteNode) (uint, error) {
	const op = "storage.postgres.CreateNoteNode"

	var nodesCount uint
	err := s.db.Get(&nodesCount, getNodesCountByIdQuery, newNoteNode.NoteId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id uint

	err = s.db.Get(&id, createBlankNoteNodeQuery, newNoteNode.NoteId, nodesCount, newNoteNode.ContentType, newNoteNode.Content)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetNotesByUserId(userId string) ([]note.Note, error) {
	const op = "storage.postgres.GetNotesByUserId"

	var notes []note.Note

	err := s.db.Select(&notes, getNotesByUserIdQuery, userId)
	if err != nil {
		return notes, fmt.Errorf("%s: %w", op, err)
	}

	return notes, nil
}

func (s *Storage) GetNoteById(id uint) (note.Note, error) {
	const op = "storage.postgres.GetNote"

	var note note.Note

	err := s.db.Get(&note, getNoteQuery, id)
	if err != nil {
		return note, fmt.Errorf("%s: %w", op, err)
	}

	return note, nil
}

func (s *Storage) UpdateFullNote(note note.Note) error {
	const op = "storage.postgres.UpdateNote"

	tx := s.db.MustBegin()

	_, err := tx.Exec(updateNoteTitleQuery, note.ID, note.Title)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, noteNode := range note.Nodes {
		_, err := tx.Exec(updateNoteNodeQuery, noteNode.Id, noteNode.ContentType, noteNode.Content)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ArchiveNote(id uint) error {
	const op = "storage.postgres.ArchiveNote"

	_, err := s.db.Exec(archiveNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteNote(id uint) error {
	const op = "storage.postgres.DeleteNote"

	_, err := s.db.Exec(deleteNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
