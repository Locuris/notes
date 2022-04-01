package test

import (
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

func TestRegisterNewUser(t *testing.T) {

	resp, err := http.Get(testUrl + "/notes")

	handleReqError(err, t)

	checkHttpOk(resp.StatusCode, t)
}
