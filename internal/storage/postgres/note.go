package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models/note"
	"main/internal/storage"
	"main/internal/storage/postgres/queries"

	"github.com/jmoiron/sqlx"
)

func (s *Storage) CreateNote(noteTitle string, userId string) (int, error) {
	const op = "storage.postgres.CreateNote"

	// creating note
	var id int

	err := s.db.Get(&id, queries.CreateNoteQuery, noteTitle, userId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// adding blank text note node
	_, err = s.AddNoteNode(id, note.ContentTypeText, "")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserNotes(userId string) ([]note.NotePreview, error) {
	const op = "storage.postgres.GetNotesByUserId"

	var notes []note.NotePreview

	// getting notes by user_id
	err := s.db.Select(&notes, queries.GetNotesByUserIdQuery, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return []note.NotePreview{}, nil
	}
	if err != nil {
		return notes, fmt.Errorf("%s: %w", op, err)
	}

	if notes == nil {
		return []note.NotePreview{}, nil
	}

	return notes, nil
}

func (s *Storage) GetNoteById(id int) (note.Note, error) {
	const op = "storage.postgres.GetNote"

	// getting note by id without nodes
	var noteFromDB note.Note

	err := s.db.Get(&noteFromDB, queries.GetNoteQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return noteFromDB, storage.ErrNoteNotFound
	}
	if err != nil {
		return noteFromDB, fmt.Errorf("%s: %w", op, err)
	}

	return noteFromDB, nil
}

func (s *Storage) GetPublicNote(id int) (note.Note, error) {
	const op = "storage.postgres.GetPublicNote"

	// getting note by id without nodes
	var noteFromDB note.Note

	err := s.db.Get(&noteFromDB, queries.GetPublicNoteQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return noteFromDB, storage.ErrNoteNotFound
	}
	if err != nil {
		return noteFromDB, fmt.Errorf("%s: %w", op, err)
	}

	return noteFromDB, nil
}

func (s *Storage) UpdateNoteTitle(id int, title string) error {
	const op = "storage.postgres.UpdateNoteTitle"

	// updating note title
	res, err := s.db.Exec(queries.UpdateNoteTitleQuery, id, title)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if note wasn't found
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

	// begin transaction
	tx := s.db.MustBegin()

	// updating note title
	res, err := tx.Exec(queries.UpdateNoteTitleQuery, id, note.Title)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// counting rows
	rows, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected += int(rows)

	// iterating over note nodes and updating content
	rows, err = updateAllNestedNodes(tx, op, note.Nodes)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected += int(rows)

	// commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return rowsAffected, nil
}

func (s *Storage) ArchiveNote(id int) error {
	const op = "storage.postgres.ArchiveNote"

	// archiving note
	res, err := s.db.Exec(queries.ArchiveNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if note wasn't found
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return storage.ErrNoteNotFound
	}

	return nil
}

func (s *Storage) UnarchiveNote(id int) error {
	const op = "storage.postgres.ArchiveNote"

	// unarchiving note
	res, err := s.db.Exec(queries.UnarchiveNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if note wasn't found
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return storage.ErrNoteNotFound
	}

	return nil
}

func (s *Storage) DeleteNote(id int) error {
	const op = "storage.postgres.DeleteNote"

	// deleting note
	res, err := s.db.Exec(queries.DeleteNoteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if note wasn't found
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrNoteNotFound
	}

	return nil
}

func (s *Storage) UpdateNoteNodeOrder(noteId int, oldOrder int, newOrder int) error {
	const op = "storage.postgres.UpdateNoteNodeOrder"

	// begin transaction
	tx := s.db.MustBegin()

	// check if note node with new_order out of bounds
	err := checkBounds(tx, op, noteId, newOrder)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if note node with old_order exists
	err = isNoteNodeExists(tx, op, noteId, oldOrder)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// update all note nodes' order between old_order and new_order
	_, err = tx.Exec(queries.UpdateOrderQuery, noteId, oldOrder, newOrder)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// set updated_at field on note
	_, err = tx.Exec(queries.SetUpdatedAtQuery, noteId)
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

func (s *Storage) MakeNotePublic(noteId int) error {
	const op = "storage.postgres.MakeNotePublic"

	// making note public
	res, err := s.db.Exec(queries.MakeNotePublicQuery, noteId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if note wasn't found
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrNoteNotFound
	}

	return nil
}

func (s *Storage) MakeNotePrivate(noteId int) error {
	const op = "storage.postgres.MakeNotePrivate"

	// making note private
	res, err := s.db.Exec(queries.MakeNotePrivateQuery, noteId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if note wasn't found
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrNoteNotFound
	}

	return nil
}

func (s *Storage) IsUserNoteOwner(userId string, noteId int) (bool, error) {
	const op = "storage.postgres.IsUserNoteOwner"

	var exists int

	err := s.db.Get(&exists, queries.IsUserNoteOwnerQuery, userId, noteId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists == 1, nil
}

func updateAllNestedNodes(tx *sqlx.Tx, op string, nodes []note.NoteNode) (int64, error) {
	var rowsAffected int64

	for _, noteNode := range nodes {
		res, err := tx.Exec(queries.UpdateNoteNodeContentQuery, noteNode.Id, noteNode.Content)
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		rowsAffected += rows
	}

	return rowsAffected, nil
}

func checkBounds(tx *sqlx.Tx, op string, noteId int, newOrder int) error {
	var count int

	err := tx.Get(&count, queries.NodesCountQuery, noteId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		return fmt.Errorf("%s: no note nodes found with note_id=%d", op, noteId)
	}
	if newOrder < 0 || newOrder >= count {
		return fmt.Errorf("%s: newOrder %d is out of range (must be between 0 and %d)", op, newOrder, count-1)
	}

	return nil
}

func isNoteNodeExists(tx *sqlx.Tx, op string, noteId int, oldOrder int) error {
	var exists int

	err := tx.Get(&exists, queries.GetNoteNodeByOrderQuery, noteId, oldOrder)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if exists == 0 {
		return fmt.Errorf("%s: no note node found with note_id=%d and order=%d", op, noteId, oldOrder)
	}

	return nil
}
