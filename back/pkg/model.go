package pkg

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Note struct {
	ID           NoteID    `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Content      string    `json:"content,omitempty"`
	CreationDate time.Time `json:"created_date,omitempty"`
}

type NoteID struct {
	uuid.UUID
}

type Storage interface {
	// List retrieves notes or an error otherwise
	List() ([]*Note, error) // todo pagination
	// Create creates a note and returns a reference to it (id) or error otherwise
	Create(title, content string) (*Note, error)
	// Retrieve retrieves a note specified by id
	Retrieve(NoteID) (*Note, error)
	// Delete atomically deletes multiple notes
	Delete([]NoteID) error
}
