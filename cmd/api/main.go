package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const webPort = ":8080"

type Config struct{}

func main() {
	//app := Config{}

	log.Printf("Starting api service on port %s\n", webPort)

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to read image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		out, err := os.Create("/tmp/uploaded-image")
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Image uploaded successfully")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s", webPort),
		Handler: nil,
	}

	error := server.ListenAndServe()
	if error != nil {
		log.Fatal(error)
	}
}
