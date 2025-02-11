package note

import (
	"time"
)

type ContentType string

const (
	ContentTypeText  = "text"
	ContentTypeImage = "image"
	ContentTypeList  = "list"
)

type Note struct {
	Id         int        `json:"id"`
	UserId     string     `json:"user_id" db:"user_id"`
	Title      string     `json:"title"`
	Nodes      []NoteNode `json:"nodes"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type NoteNode struct {
	Id          int         `json:"id"`
	NoteId      int         `json:"note_id" db:"note_id"`
	Order       int         `json:"order"`
	ContentType ContentType `json:"content_type" db:"content_type"`
	Content     string      `json:"content"`
}

type NotePreview struct {
	Id         int        `json:"id"`
	UserId     string     `json:"user_id" db:"user_id"`
	Title      string     `json:"title"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}
