package cmd

import (
	"fmt"
	"net/http"
)

func Open(w http.ResponseWriter, r *http.Request) {
	page := r.PostFormValue("page")

	url := fmt.Sprintf("http://localhost:8080/%s/", page)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
