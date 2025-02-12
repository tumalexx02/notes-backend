package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models/note"
	"main/internal/storage"
)

const (
	createNoteQuery = `
		INSERT INTO notes (title, user_id) 
		VALUES ($1, $2)
		RETURNING id;
	`
	getNoteQuery = `
		SELECT * FROM notes
		WHERE id = $1;
	`
	getNotesByUserIdQuery = `
		SELECT * FROM notes
		WHERE user_id = $1
		ORDER BY updated_at DESC;
	`
	getNoteNodesQuery = `
		SELECT * FROM note_nodes
		WHERE note_id = $1
		ORDER BY "order";
	`
	updateNoteTitleQuery = `
		UPDATE notes
		SET title = $2, updated_at = NOW()
		WHERE id = $1;
	`
	archiveNoteQuery = `
		UPDATE notes
		SET archived_at = NOW()
		WHERE id = $1;
	`
	unarchiveNoteQuery = `
		UPDATE notes
		SET archived_at = NULL
		WHERE id = $1;
	`
	deleteNoteQuery = `
		DELETE FROM notes
		WHERE id = $1;
	`
	setUpdatedAtQuery = `
		UPDATE notes
		SET updated_at = NOW()
		WHERE id = $1;
	`
)

func (s *Storage) CreateNote(noteTitle string, userId string) (int, error) {
	const op = "storage.postgres.CreateNote"

	var id int

	err := s.db.QueryRow(createNoteQuery, noteTitle, userId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.AddNoteNode(id, note.ContentTypeText, "")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserNotes(userId string) ([]note.NotePreview, error) {
	const op = "storage.postgres.GetNotesByUserId"

	var notes []note.NotePreview

	err := s.db.Select(&notes, getNotesByUserIdQuery, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return []note.NotePreview{}, nil
	}
	if err != nil {
		return notes, fmt.Errorf("%s: %w", op, err)
	}

	return notes, nil
}

func (s *Storage) GetNoteById(id int) (note.Note, error) {
	const op = "storage.postgres.GetNote"

	var noteFromDB note.Note

	err := s.db.Get(&noteFromDB, getNoteQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return noteFromDB, storage.ErrNoteNotFound
	}
	if err != nil {
		return noteFromDB, fmt.Errorf("%s: %w", op, err)
	}

	var noteNodes []note.NoteNode

	err = s.db.Select(&noteNodes, getNoteNodesQuery, noteFromDB.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return noteFromDB, storage.ErrNoteNotFound
	}
	if err != nil {
		return noteFromDB, fmt.Errorf("%s: %w", op, err)
	}

	noteFromDB.Nodes = noteNodes

	return noteFromDB, nil
}

func (s *Storage) UpdateNoteTitle(id int, title string) error {
	const op = "storage.postgres.UpdateNoteTitle"

	res, err := s.db.Exec(updateNoteTitleQuery, id, title)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return storage.ErrNoteNotFound
	}

	return nil
}

func (s *Storage) UpdateFullNote(id int, note note.Note) (int, error) {
	const op = "storage.postgres.UpdateNote"

	var rowsAffected int

	tx := s.db.MustBegin()

	res, err := tx.Exec(updateNoteTitleQuery, id, note.Title)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected += int(rows)

	for _, noteNode := range note.Nodes {
		res, err := tx.Exec(updateNoteNodeContentQuery, noteNode.Id, noteNode.Content)
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		rowsAffected += int(rows)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return rowsAffected, nil
}

func (s *Storage) ArchiveNote(id int) error {
	const op = "storage.postgres.ArchiveNote"

	_, err := s.db.Exec(archiveNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UnarchiveNote(id int) error {
	const op = "storage.postgres.ArchiveNote"

	_, err := s.db.Exec(archiveNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteNote(id int) error {
	const op = "storage.postgres.DeleteNote"

	res, err := s.db.Exec(deleteNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrNoteNotFound
	}

	return nil
}
