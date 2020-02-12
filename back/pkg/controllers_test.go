package pkg_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"back/pkg"

	uuid "github.com/satori/go.uuid"
)

func TestController_Create(t *testing.T) {
	api := pkg.NewController(*pkg.NewInMemoryStorage())

	_, response, err := create(api, "", "")
	if err != nil {
		t.Errorf("failed creating: %s", err)
	}

	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
	}
}

func TestController_List(t *testing.T) {
	api := pkg.NewController(*pkg.NewInMemoryStorage())

	id1, _, _ := create(api, "title1", "content1")
	id2, _, _ := create(api, "title2", "content2")

	req, _ := http.NewRequest("", "", strings.NewReader(""))
	response := doRequest(req, api.List)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var data []*pkg.Note

	decoder := json.NewDecoder(response.Body)

	err := decoder.Decode(&data)
	if err != nil && err != io.EOF {
		t.Error("errored while decoding response data")
	}

	if len(data) != 2 {
		t.Error("wrong list length")
	}

	// list can come in weird order (since the underlying data structure is a map), so we need this kind of matching:
	for _, createdID := range []pkg.NoteID{{*id1}, {*id2}} {
		match := false
		for _, listNote := range data {
			if createdID == listNote.ID {
				match = true
			}
		}
		if match == false {
			t.Errorf("a freshly created id (%s) is missing from list", createdID)
		}
	}
}

func TestController_Delete(t *testing.T) {
	api := pkg.NewController(*pkg.NewInMemoryStorage())

	existingID, _, err := create(api, "", "")
	if err != nil {
		t.Error("failed creating note")
	}

	req, _ := http.NewRequest("", "", strings.NewReader(fmt.Sprintf(`{"ids": ["%s"]}`, existingID)))
	response := doRequest(req, api.Delete)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func create(a *pkg.Controller, title, content string) (id *uuid.UUID, response *http.Response, err error) {
	req, _ := http.NewRequest("", "", strings.NewReader(
		fmt.Sprintf(`{"title": "%s", "content": "%s"}`, title, content)))

	response = doRequest(req, a.Create)

	var data struct{ ID uuid.UUID }

	body, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, err
	}

	return &data.ID, response, nil
}

func doRequest(req *http.Request, handler http.HandlerFunc) *http.Response {
	w := httptest.NewRecorder()
	handler(w, req)
	return w.Result()
}
