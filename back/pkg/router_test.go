package pkg_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "back/pkg"

	uuid "github.com/satori/go.uuid"
)

type testServer struct {
	*httptest.Server
	// state for the last performed operation:
	createdIDs []uuid.UUID
	listNotes  []*Note
}

// TestController_E2E will spawn a server on the local interface on random ports
func TestController_E2E(t *testing.T) {
	api := NewController(*NewInMemoryStorage())
	s := &testServer{Server: httptest.NewServer(api.Router)}
	defer s.Close()

	t.Run("empty list", s.list(0))
	t.Run("create 1", s.create([]Note{{Title: "", Content: ""}}))
	t.Run("check 1 item", s.list(1))

	if s.listNotes[0].ID.UUID != s.createdIDs[0] {
		t.Errorf("created ID not found in the list")
	}

	t.Run("delete the item", s.delete(s.listNotes[0].ID, true))
	t.Run("list should be empty", s.list(0))
	t.Run("create 1 more", s.create([]Note{{Title: "nothing", Content: ""}}))
	t.Run("delete it using body", s.delete(NoteID{UUID: s.createdIDs[0]}, false))

	t.Run("list is empty again", s.list(0))
}

func (s *testServer) list(count int) func(t *testing.T) {
	return func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, s.URL, nil)
		if err != nil {
			t.Fatalf("can't create request: %s", err)
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("can't perform request: %s", err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		if err = json.NewDecoder(res.Body).Decode(&s.listNotes); err != nil {
			t.Errorf("error unmarshaling response json")
		}

		if len(s.listNotes) != count {
			t.Errorf("list count is wrong (expected %d, got %d items)", count, len(s.listNotes))
		}
	}
}

func (s *testServer) create(notes []Note) func(t *testing.T) {
	return func(t *testing.T) {
		s.createdIDs = nil

		for _, note := range notes {
			res, err := http.Post(s.URL, "application/json", strings.NewReader(
				fmt.Sprintf(`{"title":"%s", "content":"%s"}`, note.Title, note.Content)))
			if err != nil {
				log.Fatalf("can't POST create request: %s", err)
			}

			if res.StatusCode != http.StatusCreated {
				t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
			}

			var incoming struct{ Id uuid.UUID }
			if err = json.NewDecoder(res.Body).Decode(&incoming); err != nil {
				t.Error("error unmarshaling response json")
			}

			s.createdIDs = append(s.createdIDs, incoming.Id)
		}
	}
}

func (s *testServer) delete(id NoteID, inputInQuery bool) func(t *testing.T) {
	return func(t *testing.T) {
		var req *http.Request
		var err error

		if inputInQuery {
			targetURL := fmt.Sprintf("%s?ids=%s", s.URL, id)
			req, err = http.NewRequest(http.MethodDelete, targetURL, nil)
		} else {
			body := fmt.Sprintf(`{"ids": ["%s"]}`, id)
			req, err = http.NewRequest(http.MethodDelete, s.URL, strings.NewReader(body))
		}

		if err != nil {
			t.Fatalf("can't create request: %s", err)
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("error performing request: %s", err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}
	}
}
