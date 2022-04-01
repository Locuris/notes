package main

import (
	"log"
	c "notes/config"
	"notes/server"
)

func main() {

	config, err := c.GetConfig()

	if err != nil {
		log.Panicln("cannot load config", err)
	}

	port := config.Server.Port

	server.StartHttpServer(port)
}
