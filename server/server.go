package server

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"notes/notes"
)

type NoteMessage struct {
	User notes.User `json:"user"`
	Note notes.Note `json:"note"`
}

func StartHttpServer(port string) {
	echoServer := echo.New()

	echoServer.Use(middleware.Logger())
	echoServer.GET("/health", healthCheck)

	echoServer.GET("/users/new", registerUserHandler)

	echoServer.GET("/notes/new", createNoteHandler)

	echoServer.PUT("/notes", updateNoteHandler)

	echoServer.GET("/notes/:noteId", getNoteHandler)

	echoServer.Logger.Fatal(echoServer.Start(":" + port))
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}

func registerUserHandler(c echo.Context) error {
	return c.String(http.StatusOK, notes.RegisterNewUser())
}

func createNoteHandler(c echo.Context) error {
	notepad := getNotepad(c)

	note := notepad.CreateNewNote("", "")

	return c.String(http.StatusOK, note.Id)
}

func updateNoteHandler(c echo.Context) error {
	noteMessage := getNoteMessage(c)

	notepad := notes.GetNotepad(noteMessage.User.Id)
	note, err := notepad.GetNote(noteMessage.Note.Id)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	note.Title = noteMessage.Note.Title
	note.Content = noteMessage.Note.Content
	note.Archived = noteMessage.Note.Archived

	return c.NoContent(http.StatusOK)
}

func getNoteHandler(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusNotFound)
	}

	notepad := notes.GetNotepad(user.Id)
	if notepad == nil {
		return c.NoContent(http.StatusNotFound)
	}

	noteId := getParam(c, "noteId")

	note, err := notepad.GetNote(noteId)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, note)
}

func getQueryParam(c echo.Context, paramName string) string {
	param := c.QueryParam(paramName)
	log.Println("param is:", param)
	return param
}

func getParam(c echo.Context, paramName string) string {
	param := c.Param(paramName)
	log.Println("param is:", param)
	return param
}

func getNotepad(c echo.Context) *notes.Notepad {
	userId := getQueryParam(c, "userId")
	if userId == "" {
		return nil
	}

	return notes.GetNotepad(userId)
}

func getNoteFromId(c echo.Context, notepad *notes.Notepad) *notes.Note {
	noteId := getParam(c, "noteId")

	note, err2 := notepad.GetNote(noteId)
	if err2 != nil {
		return nil
	}
	return note
}

func getNoteMessage(c echo.Context) *NoteMessage {
	bodyBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil
	}

	noteMessage := new(NoteMessage)

	err = json.Unmarshal(bodyBytes, noteMessage)
	if err != nil {
		return nil
	}

	return noteMessage
}

func getUser(c echo.Context) *notes.User {
	bodyBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil
	}

	user := new(notes.User)

	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		return nil
	}

	return user
}
