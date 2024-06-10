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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("dist/static"))))

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
