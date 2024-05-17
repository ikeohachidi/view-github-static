package cmd

import (
	"fmt"
	"net/http"
	"os"
)

func Open(w http.ResponseWriter, r *http.Request) {
	page := r.PostFormValue("page")

	urlEnv := os.Getenv("URL")

	url := fmt.Sprintf("%s/%s/", urlEnv, page)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
