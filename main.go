package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ikeohachidi/view-github-static/cmd"
)

func main() {
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./web"))

	router.Handle("GET /", fs)

	router.HandleFunc("POST /open", cmd.Open)
	router.HandleFunc("GET /{username}/{repo}/*", cmd.Navigate)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Print("Listening on port 8080")
	log.Fatal(server.ListenAndServe())
}
