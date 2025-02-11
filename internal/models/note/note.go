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
	ID         int        `json:"id"`
	UserId     string     `json:"user_id"`
	Title      string     `json:"title"`
	Nodes      []NoteNode `json:"nodes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty"`
}

type NoteNode struct {
	Id          int         `json:"id"`
	NoteId      int         `json:"note_id"`
	Order       int         `json:"order"`
	ContentType ContentType `json:"type"`
	Content     string      `json:"content"`
}
