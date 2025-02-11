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

// GetUserNotes gets all notes for a given user ID from the database. If the user doesn't have any notes,
// it returns an empty slice and ErrNoteNotFound. If the operation fails, it returns an error.

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

// GetNoteById gets a note with the given ID from the database. If the note with the given ID
// doesn't exist, it returns ErrNoteNotFound. If the operation fails, it returns an error.

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

// UpdateNoteTitle updates the title of a note with the given ID in the database.
// If the update operation fails, it returns an error.

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

// UpdateFullNote updates a note with a given ID to the given Note object.
// It is atomic, meaning that either all of the updates are committed or none of them are.
// If any part of the update fails, the entire transaction is rolled back and
// an error is returned.

func (s *Storage) UpdateFullNote(id int, note note.Note) error {
	const op = "storage.postgres.UpdateNote"

	tx := s.db.MustBegin()

	_, err := tx.Exec(updateNoteTitleQuery, note.Id, note.Title)
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

// ArchiveNote sets the archived_at field of a note to the current time for the given
// note ID, effectively archiving the note. Returns an error if the operation fails.

func (s *Storage) ArchiveNote(id int) error {
	const op = "storage.postgres.ArchiveNote"

	_, err := s.db.Exec(archiveNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// UnarchiveNote sets the archived_at field of a note to NULL for the given note ID,
// effectively unarchiving the note. Returns an error if the operation fails.

func (s *Storage) UnarchiveNote(id int) error {
	const op = "storage.postgres.ArchiveNote"

	_, err := s.db.Exec(archiveNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// DeleteNote deletes the note with the given ID from the database. If the note
// doesn't exist, it returns ErrNoteNotFound.

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
