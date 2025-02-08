package notes

import (
	"time"
)

type ContentType string

var (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
	ContentTypeList  ContentType = "list"
)

type Note struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title"`
	Nodes      []NoteNode `json:"nodes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty"`
}

type NoteNode struct {
	Id          uint        `json:"id"`
	Order       uint        `json:"order"`
	ContentType ContentType `json:"type"`
	Content     string      `json:"content"`
}
