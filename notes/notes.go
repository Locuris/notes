package notes

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

type Notepad struct {
	User  *User
	Notes map[string]*Note
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Note struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Archived bool   `json:"archived"`
}

var NoteNotFoundError = errors.New("note not found")

var notepads = make(map[string]*Notepad)

func RegisterNewUser() string {
	user := createNewUser()
	notepad := user.createNewNotepad()
	notepads[notepad.User.Id] = &notepad
	return notepad.User.Id
}

func GetNotepad(userId string) *Notepad {
	notepad, ok := notepads[userId]
	if !ok {
		log.Panicf("no notepad found with userId: %q", userId)
		return nil
	}
	return notepad
}

func (user *User) createNewNotepad() Notepad {
	return Notepad{
		User:  user,
		Notes: make(map[string]*Note),
	}
}

func createNewUser() *User {
	return &User{
		Id:   uuid.New().String(),
		Name: "",
	}
}

func (notepad *Notepad) CreateNewNote(title string, content string) *Note {
	newId := uuid.New().String()
	note := &Note{
		Id:       newId,
		Title:    title,
		Content:  content,
		Archived: false,
	}
	notepad.Notes[newId] = note
	return note
}

func (notepad *Notepad) GetNote(id string) (*Note, error) {
	note, ok := notepad.Notes[id]
	if !ok {
		log.Panicf("No note with id %q found", id)
		return note, NoteNotFoundError
	}
	return note, nil
}

func (notepad *Notepad) UpdateNote(id string, title string, content string) error {
	note, err := notepad.GetNote(id)
	if err != nil {
		return NoteNotFoundError
	}
	note.Title = title
	note.Content = content
	return nil
}

func (notepad *Notepad) setArchiveState(id string, archive bool) error {
	note, err := notepad.GetNote(id)
	if err != nil {
		return NoteNotFoundError
	}
	if note.Archived == archive {
		log.Printf("unneccesary attempt to chang archived archive of note %(. this should not happen")
	}
	note.Archived = archive
	return nil
}

func (notepad *Notepad) ArchiveNote(id string) error {
	return notepad.setArchiveState(id, true)
}

func (notepad Notepad) UnarchiveNote(id string) error {
	return notepad.setArchiveState(id, false)
}

func (notepad *Notepad) getNotesByArchiveState(isArchived bool) []*Note {
	notes := make([]*Note, 0)
	for _, note := range notepad.Notes {
		if note.Archived == isArchived {
			notes = append(notes, note)
		}
	}
	return notes
}

func (notepad *Notepad) GetSavedNotes() []*Note {
	return notepad.getNotesByArchiveState(false)
}

func (notepad *Notepad) GetArchivedNotes() []*Note {
	return notepad.getNotesByArchiveState(true)
}

func (notepad *Notepad) DeleteNote(id string) error {
	_, err := notepad.GetNote(id)
	if err != nil {
		return err
	}
	delete(notepad.Notes, id)
	return nil
}
