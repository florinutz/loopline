package pkg_test

import (
	"reflect"
	"testing"

	"back/pkg"
)

// unit tests for the in-mem implementation

func TestService_Create_Retrieve_Delete(t *testing.T) {
	svc := pkg.NewInMemoryStorage()

	note, err := svc.Create("title", "content")
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}

	note, err = svc.Retrieve(note.ID)
	if err != nil {
		t.Fatalf("Retrieve failed: %s", err)
	}

	if err = svc.Delete([]pkg.NoteID{note.ID}); err != nil {
		t.Fatalf("Got error while Deleting")
	}

	note, err = svc.Retrieve(note.ID) // negative path
	if err != pkg.ErrNoteDoesntExist {
		t.Fatalf("Err is not ErrNoteDoesnExist: %s", err)
	}
}

func TestService_List(t *testing.T) {
	svc := pkg.NewInMemoryStorage()

	note1, _ := svc.Create("title1", "content1")
	note2, _ := svc.Create("title2", "content2")

	list, err := svc.List()
	if err != nil {
		t.Fatalf("List failed: %s", err)
	}

	// this also tests the order, which should be desc by time
	if !reflect.DeepEqual([]*pkg.Note{note2, note1}, list) {
		t.Error("there's something wrong with the list")
	}
}
