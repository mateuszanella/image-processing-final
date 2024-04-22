package main

import (
	"log"
	"net/http"
)

const webPort = ":8080"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting api service on port %s\n", webPort)

	server := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
