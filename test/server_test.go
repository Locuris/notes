package test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
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

func getNoteId(t *testing.T) (noteId string, userId string) {
	testUserId := getUserId(t)

	resp, err := http.Get(testUrl + "/notes/new?userId=" + testUserId)

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)

	return responseBodyToString(resp, t), testUserId
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

	testNoteId, _ := getNoteId(t)

	t.Logf(testNoteId)

	testUuidFormat(testNoteId, t)
}

func TestUpdateNoteHandler(t *testing.T) {
	testNoteId, testUserId := getNoteId(t)

	testNoteMessage := server.NoteMessage{
		User: notes.User{
			Id:   testUserId,
			Name: "",
		},
		Note: notes.Note{
			Id:       testNoteId,
			Title:    "New Title",
			Content:  "New Content",
			Archived: false,
		},
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
}

func TestArchiveNoteHandler(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestUnarchiveNoteHandler(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestGetSavedNotesHandler(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestGetArchivedNotesHandler(t *testing.T) {
	t.Error("Not implemented yet")
}
