package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Controller struct {
	Storage
	Router *mux.Router
}

func NewController(notebookService InMemoryStorage) *Controller {
	c := &Controller{Storage: &notebookService}
	c.Router = c.generateRouter()
	return c
}

func (c *Controller) Create(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		respondError(res, http.StatusBadRequest, "missing body")
		return
	}

	var inputData struct{ Title, Content string }

	err := json.NewDecoder(req.Body).Decode(&inputData)
	if err != nil && err != io.EOF {
		respondError(res, http.StatusBadRequest, "the request "+
			"should be a valid json containing the *title* and *body* keys")
		return
	}

	note, err := c.Storage.Create(inputData.Title, inputData.Content)
	if err != nil {
		respondError(res, http.StatusInternalServerError, "error creating note")
		return
	}

	if err = respond(res, http.StatusCreated, map[string]NoteID{"id": note.ID}); err != nil {
		fmt.Fprintf(os.Stderr, "error responding with the created note: %s", err)
	}
}

func (c *Controller) List(res http.ResponseWriter, req *http.Request) {
	notes, err := c.Storage.List()
	if err != nil {
		respondError(res, http.StatusInternalServerError, "error listing notes")
		return
	}

	if err = respond(res, http.StatusOK, notes); err != nil {
		fmt.Fprintf(os.Stderr, "error responding with the list of notes: %s", err)
	}
}

func (c *Controller) Delete(res http.ResponseWriter, req *http.Request) {
	// if the array is present in the url, use it then stop
	ids, ok := req.URL.Query()["ids"]
	if ok && len(ids) >= 1 {
		var properIDs []NoteID

		for _, id := range ids {
			uu, err := uuid.FromString(id)
			if err != nil {
				// one of the ids was invalid. stop.
				respondError(res, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", id))
				return
			}
			properIDs = append(properIDs, NoteID{uu})
		}

		if err := c.Storage.Delete(properIDs); err != nil {
			respondError(res, http.StatusInternalServerError, "error deleting notes")
			return
		}

		return
	}

	// the *ids* key was not present in the url, so try to grab the data from req body
	if req.Body == nil {
		respondError(res, http.StatusBadRequest, "missing body")
		return
	}

	var inputData struct {
		IDs []NoteID `json:"ids"`
	}

	err := json.NewDecoder(req.Body).Decode(&inputData)
	if err != nil {
		respondError(res, http.StatusBadRequest, "can't handle incoming request")
		return
	}

	if err = c.Storage.Delete(inputData.IDs); err != nil {
		respondError(res, http.StatusInternalServerError, "error deleting notes")
		return
	}

	if err = respond(res, http.StatusOK, nil); err != nil {
		fmt.Fprintf(os.Stderr, "error responding deletion request: %s", err)
	}
}

func respondError(w http.ResponseWriter, code int, message string) {
	_ = respond(w, code, map[string]string{"error": message})
}

// todo make sure this doesn't leak implementation details in production
func respond(rw http.ResponseWriter, code int, payload interface{}) error {
	if payload == nil {
		rw.WriteHeader(code)
		return nil
	}

	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	_, _ = rw.Write(response)

	return nil
}
