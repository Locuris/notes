package test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"notes/notes"
	"notes/server"
	"testing"
	"time"
)

const testUrl = "http://localhost:8080"

func startServer() {
	go func() {
		server.StartHttpServer("8080")
	}()

	time.Sleep(1 * time.Second)
}

func handleReqError(err error, t *testing.T) {
	if err != nil {
		t.Error("could not create request", err)
	}
}

func responseBodyToString(resp *http.Response, t *testing.T) string {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error("could note parse response body", err)
	}
	return string(b)
}

func checkHttpOk(statusCode int, t *testing.T) {
	if statusCode != http.StatusOK {
		t.Error("response http status not ok, is", statusCode)
	}
}

func testUuidFormat(uuidToTest string, t *testing.T) {
	testUuid := uuid.New().String()

	l1 := len(uuidToTest)
	l2 := len(testUuid)

	t.Log("comparing lengths", l1, l2)

	if len(uuidToTest) != len(testUuid) {
		t.Errorf("format of user id returned %q does not match standard uuid format e.g: %q", uuidToTest, testUuid)
	}
}

func getUserId(t *testing.T) string {
	resp, err := http.Get(testUrl + "/users/new")

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)

	return responseBodyToString(resp, t)
}

func getNoteAndUserIds(t *testing.T) (noteId string, userId string) {
	testUserId := getUserId(t)

	resp, err := http.Get(testUrl + "/notes/new?userId=" + testUserId)

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)

	return responseBodyToString(resp, t), testUserId
}

func getNoteId(t *testing.T, userId string) (noteId string) {
	resp, err := http.Get(testUrl + "/notes/new?userId=" + userId)

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)

	return responseBodyToString(resp, t)
}

func updateNote(t *testing.T, userId string, noteId string, title string, content string, archived bool) (notes.User, notes.Note) {
	user := notes.User{
		Id:   userId,
		Name: "",
	}

	note := notes.Note{
		Id:       noteId,
		Title:    title,
		Content:  content,
		Archived: archived,
	}

	testNoteMessage := server.NoteMessage{
		User: user,
		Note: note,
	}

	js, err := json.Marshal(testNoteMessage)
	if err != nil {
		t.Error("Could note marshal json", err)
	}

	req, err1 := http.NewRequest(http.MethodPut, testUrl+"/notes", bytes.NewBuffer(js))
	if err1 != nil {
		t.Error("Could not create the request", err1)
	}

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		t.Error("Could not send request", err2)
	}

	checkHttpOk(resp.StatusCode, t)

	return user, note
}

func TestStartHttpServer(t *testing.T) {
	startServer()

	resp, err := http.Get(testUrl + "/health")

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)
}

func TestRegisterUserHandler(t *testing.T) {

	testUserId := getUserId(t)

	t.Logf(testUserId)

	testUuidFormat(testUserId, t)

}

func TestCreateNewNoteHandler(t *testing.T) {

	testNoteId, _ := getNoteAndUserIds(t)

	t.Logf(testNoteId)

	testUuidFormat(testNoteId, t)
}

func TestUpdateNoteHandler(t *testing.T) {
	testNoteId, testUserId := getNoteAndUserIds(t)

	user, testNote := updateNote(t, testUserId, testNoteId, "New Title", "New Content", true)

	notepad := notes.GetNotepad(user.Id)
	note, err3 := notepad.GetNote(testNote.Id)
	if err3 != nil {
		t.Error("Could note get note", err3)
	}

	if note.Title != testNote.Title || note.Content != testNote.Content || !note.Archived {
		t.Errorf("Note title or content does not match. %q should be %q and %q should be %q", note.Title, testNote.Title, note.Content, testNote.Content)
	}
}

func TestGetSavedNotesHandler(t *testing.T) {
	testNote1Id, testUserId := getNoteAndUserIds(t)
	testNote2Id := getNoteId(t, testUserId)
	testNote3Id := getNoteId(t, testUserId)

	user, _ := updateNote(t, testUserId, testNote1Id, "Note 1", "This is note 1 content.", false)
	updateNote(t, testUserId, testNote2Id, "Note 2", "This is note 2 content.", true)
	updateNote(t, testUserId, testNote3Id, "Note 3", "This is note 3 content.", false)

	js, err := json.Marshal(user)
	if err != nil {
		t.Error("Could note marshal json", err)
	}

	req, err1 := http.NewRequest(http.MethodGet, testUrl+"/notes/saved", bytes.NewBuffer(js))
	if err1 != nil {
		t.Error("Could not create the request", err1)
	}

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		t.Error("Could not send request", err2)
	}

	checkHttpOk(resp.StatusCode, t)

	bodyBytes, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		t.Error("could not read response body", err3)
	}

	testNotes := new(server.NoteCollection)

	err = json.Unmarshal(bodyBytes, testNotes)
	if err != nil {
		t.Error("could not parse json in response", err)
	}

	notesLength := len(testNotes.Notes)
	if notesLength != 2 {
		t.Errorf("incorrect number of notes returned, expecting %q but got %q", 2, notesLength)
	}

	for _, note := range testNotes.Notes {
		if note.Archived {
			t.Error("note is archived, it should not be")
		}
	}

	t.Log(testNotes)
}

func TestGetArchivedNotesHandler(t *testing.T) {
	testNote1Id, testUserId := getNoteAndUserIds(t)
	testNote2Id := getNoteId(t, testUserId)
	testNote3Id := getNoteId(t, testUserId)

	user, _ := updateNote(t, testUserId, testNote1Id, "Note 1", "This is note 1 content.", false)
	updateNote(t, testUserId, testNote2Id, "Note 2", "This is note 2 content.", true)
	updateNote(t, testUserId, testNote3Id, "Note 3", "This is note 3 content.", false)

	js, err := json.Marshal(user)
	if err != nil {
		t.Error("Could note marshal json", err)
	}

	req, err1 := http.NewRequest(http.MethodGet, testUrl+"/notes/archived", bytes.NewBuffer(js))
	if err1 != nil {
		t.Error("Could not create the request", err1)
	}

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		t.Error("Could not send request", err2)
	}

	checkHttpOk(resp.StatusCode, t)

	bodyBytes, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		t.Error("could not read response body", err3)
	}

	testNotes := new(server.NoteCollection)

	err = json.Unmarshal(bodyBytes, testNotes)
	if err != nil {
		t.Error("could not parse json in response", err)
	}

	notesLength := len(testNotes.Notes)
	if notesLength != 1 {
		t.Errorf("incorrect number of notes returned, expecting %q but got %q", 1, notesLength)
	}

	for _, note := range testNotes.Notes {
		if !note.Archived {
			t.Error("note is not archived, it should be")
		}
	}

	t.Log(testNotes)
}
