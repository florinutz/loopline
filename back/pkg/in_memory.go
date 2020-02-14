package pkg

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// InMemoryStorage implements the Storage and stores data in memory
type InMemoryStorage struct {
	sync.RWMutex
	notes map[NoteID]*Note
}

// NewInMemoryStorage is the constructor for InMemoryStorage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		notes: map[NoteID]*Note{},
	}
}

func (s *InMemoryStorage) Create(title, content string) (note *Note, err error) {
	s.Lock()
	defer s.Unlock()

	note = &Note{
		ID:           NoteID{uuid.NewV4()},
		Title:        title,
		Content:      content,
		CreationDate: time.Now(),
	}

	s.notes[note.ID] = note

	return
}

func (s *InMemoryStorage) List() ([]*Note, error) {
	s.Lock()
	defer s.Unlock()

	var result []*Note

	for _, note := range s.notes {
		result = append(result, note)
	}

	// order by date desc
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreationDate.After(result[j].CreationDate)
	})

	return result, nil
}

func (s *InMemoryStorage) Retrieve(ID NoteID) (*Note, error) {
	s.Lock()
	defer s.Unlock()

	note, exists := s.notes[ID]
	if !exists {
		return nil, ErrNoteDoesntExist
	}

	return note, nil
}

func (s *InMemoryStorage) Delete(IDs []NoteID) error {
	s.Lock()
	defer s.Unlock()

	var missingIds []NoteID

	for _, id := range IDs {
		if _, idExists := s.notes[id]; !idExists {
			missingIds = append(missingIds, id)
			continue
		}
		delete(s.notes, id)
	}

	var err error

	if len(missingIds) > 0 {
		var ids []string
		for _, id := range missingIds {
			ids = append(ids, id.String())
		}
		err = errors.New(fmt.Sprintf("missing ids: %s", strings.Join(ids, ", ")))
	}

	return err
}

var ErrNoteDoesntExist = errors.New("note doesn't exist")
