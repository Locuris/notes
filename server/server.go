package server

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"notes/notes"
)

func StartHttpServer(port string) {
	echoServer := echo.New()

	echoServer.Use(middleware.Logger())
	echoServer.GET("/health", healthCheck)

	echoServer.GET("/notes", registerUserHandler)

	echoServer.POST("/notes/:userId", createNoteHandler)

	echoServer.Logger.Fatal(echoServer.Start(":" + port))
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}

func registerUserHandler(c echo.Context) error {
	return c.String(http.StatusOK, notes.RegisterNewUser().String())
}

func createNoteHandler(c echo.Context) error {
	userId := getParams(c, "userId")

	userIdAsUuid, err := uuid.FromBytes([]byte(userId))

	if err != nil {
		return err
	}

	notepad := notes.GetNotepad(userIdAsUuid)

	//title := getParams(c, "title")
	//content := getParams(c, "content")

	title := ""
	content := ""

	note := notepad.CreateNewNote(title, content)

	return c.String(http.StatusOK, note.Id)
}

func getParams(c echo.Context, paramName string) string {
	param := c.QueryParam(paramName)
	log.Println("param is:", param)
	return param
}
