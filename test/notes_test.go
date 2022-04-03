package test

import (
	"notes/notes"
	"strconv"
	"testing"
)

var testNote1 = notes.Note{
	Id:       "",
	Title:    "Title 1",
	Content:  "Example content 1...",
	Archived: false,
}

var testNote2 = notes.Note{
	Id:       "",
	Title:    "Title 2",
	Content:  "Example content 2...",
	Archived: true,
}

var testNote3 = notes.Note{
	Id:       "",
	Title:    "New Title",
	Content:  "New Content",
	Archived: false,
}

func getNotepad(t *testing.T) notes.Notepad {
	id := notes.RegisterNewUser()
	notepad := notes.GetNotepad(id)
	if notepad == nil {
		t.Errorf("Could not fetch notepad")
	}
	return *notepad
}

func TestCreateNewNotepad(t *testing.T) {
	notepad := getNotepad(t)

	if notepad.Notes == nil || len(notepad.Notes) != 0 || notepad.User.Id == "" {
		t.Errorf("Notepad wasn't created properly")
	}
}

func TestCreateNewNote(t *testing.T) {
	notepad := getNotepad(t)

	note := notepad.CreateNewNote(testNote1.Title, testNote1.Content)

	if note == nil {
		t.Errorf("Note failed to instantiate")
	}

	if note.Archived || note.Title != testNote1.Title || note.Content != testNote1.Content {
		t.Errorf("Note data not what is expected. %q should be %q and %q should be %q", testNote1.Title, note.Title, testNote1.Content, note.Content)
	}
}

func TestChangeNote(t *testing.T) {
	notepad := getNotepad(t)

	note := notepad.CreateNewNote(testNote1.Title, testNote1.Content)
	id := note.Id

	err := notepad.UpdateNote(id, testNote3.Title, testNote3.Content)

	if err != nil {
		t.Error("error raised when trying to update note", err)
	}

	if note.Title != testNote3.Title || note.Content != testNote3.Content {
		t.Errorf("note has not been updated. %q should be %q and %q should be %q", note.Title, testNote3.Title, note.Content, testNote3.Content)
	}
}

func TestArchive(t *testing.T) {
	notepad := getNotepad(t)

	note := notepad.CreateNewNote(testNote1.Title, testNote1.Content)
	id := note.Id

	err := notepad.ArchiveNote(id)

	if err != nil {
		t.Error("error raised when trying to archive note", err)
	}

	if !note.Archived {
		t.Errorf("note should be archived, it is not")
	}
}

func TestUnarchive(t *testing.T) {
	notepad := getNotepad(t)

	note := notepad.CreateNewNote(testNote1.Title, testNote1.Content)
	id := note.Id

	err := notepad.ArchiveNote(id)

	err2 := notepad.UnarchiveNote(id)

	if err != nil || err2 != nil {
		t.Error("error(s) raised when performing archive state change", err, err2)
	}

	if note.Archived {
		t.Errorf("note should not be archived, it is")
	}
}

func TestGetSavedMessages(t *testing.T) {
	notepad := getNotepad(t)

	notepad.CreateNewNote(testNote1.Title, testNote1.Content)
	notepad.CreateNewNote(testNote2.Title, testNote2.Content)
	note := notepad.CreateNewNote(testNote3.Title, testNote3.Content)
	id := note.Id

	err := notepad.ArchiveNote(id)

	if err != nil {
		t.Error("error raised when trying to archive note", err)
	}

	savedNotes := notepad.GetSavedNotes()
	l := len(savedNotes)

	if l != 2 {
		t.Errorf("incorrect amount of savedNotes returned was %s should be '2'", strconv.Itoa(l))
	}

	for _, n := range savedNotes {
		if n.Archived {
			t.Errorf("archived note returned from GetSavedNotes")
		}
	}
}

func TestGetArchivedMessages(t *testing.T) {
	notepad := getNotepad(t)

	notepad.CreateNewNote(testNote1.Title, testNote1.Content)
	notepad.CreateNewNote(testNote2.Title, testNote2.Content)
	note := notepad.CreateNewNote(testNote3.Title, testNote3.Content)
	id := note.Id

	err := notepad.ArchiveNote(id)

	if err != nil {
		t.Error("error raised when trying to archive note", err)
	}

	archivedNotes := notepad.GetArchivedNotes()
	l := len(archivedNotes)

	if l != 1 {
		t.Errorf("incorrect amount of archivedNotes returned was %s should be '1'", strconv.Itoa(l))
	}

	for _, n := range archivedNotes {
		if !n.Archived {
			t.Errorf("unarchived note returned from GetArchivedNotes")
		}
	}
}
