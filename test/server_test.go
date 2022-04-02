package test

import (
	"github.com/google/uuid"
	"io"
	"net/http"
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
		t.Error("health check failed")
	}
}

func TestStartHttpServer(t *testing.T) {
	startServer()

	resp, err := http.Get(testUrl + "/health")

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)
}

func TestRegisterUserHandler(t *testing.T) {

	resp, err := http.Get(testUrl + "/notes")

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)

	userIdResp := responseBodyToString(resp, t)

	testUuid := uuid.New().String()

	if len(userIdResp) != len(testUuid) {
		t.Errorf("format of user id returned %q does not match standard uuid format e.g: %q", userIdResp, testUuid)
	}
}

func TestCreateNewNoteHandler(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestUpdateNoteHandler(t *testing.T) {
	t.Error("Not implemented yet")
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
