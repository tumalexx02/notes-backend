package note

import (
	"time"
)

type ContentType string

const (
	ContentTypeText  = "text"
	ContentTypeImage = "image"
)

type Note struct {
	Id         int        `json:"id"`
	PublicId   string     `json:"public_id" db:"public_id"`
	UserId     string     `json:"user_id" db:"user_id"`
	Title      string     `json:"title" validate:"max=31"`
	Nodes      []NoteNode `json:"nodes"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

type NoteNode struct {
	Id          int         `json:"id"`
	NoteId      int         `json:"note_id" db:"note_id"`
	Order       int         `json:"order" validate:"gte=0"`
	ContentType ContentType `json:"content_type" db:"content_type"`
	Content     string      `json:"content,omitempty"`
	Image       string      `json:"image,omitempty"`
}

type NotePreview struct {
	Id         int        `json:"id"`
	UserId     string     `json:"user_id" db:"user_id"`
	Title      string     `json:"title" validate:"max=31"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}
