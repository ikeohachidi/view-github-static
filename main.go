package main

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"embed"

	"github.com/ikeohachidi/view-github-static/cmd"
)

//go:embed "web"
var web embed.FS

func main() {
	router := http.NewServeMux()
	fSys, err := fs.Sub(web, "web")
	if err != nil {
		log.Fatalf("unable to read embed files: %v", err)
	}

	router.Handle("GET /", http.FileServer(http.FS(fSys)))

	router.HandleFunc("POST /open", cmd.Open)
	router.HandleFunc("GET /{username}/{repo}/*", cmd.Navigate)

	server := http.Server{
		Addr:         ":8000",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Print("Listening on port 8000")
	log.Fatal(server.ListenAndServe())
}
