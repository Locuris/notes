package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"notes/notes"
)

func StartHttpServer(port string) {
	echoServer := echo.New()

	echoServer.Use(middleware.Logger())

	echoServer.GET("/health", healthCheck)

	echoServer.GET("/notes", registerUserHandler)

	echoServer.POST("/notes/:userId:title:content", createNoteHandler)

	echoServer.Logger.Fatal(echoServer.Start(":" + port))
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "healthy")
}

func registerUserHandler(c echo.Context) error {
	return c.String(http.StatusOK, notes.RegisterNewUser().String())
}

func createNoteHandler(c echo.Context) error {

	return c.NoContent(http.StatusOK)
}
