package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const GH_API_BASE = "https://api.github.com"

type FileMetaData struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

func Navigate(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	repo := r.PathValue("repo")

	if len(username) == 0 || len(repo) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the rest of the path ignoring the parts with the username and repo
	// which would be the first two elements in the url path when split
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")[2:]
	pathString := strings.Join(path, "/")

	if filepath.Ext(strings.Join(path, "/")) == "" {
		pathString = strings.TrimSuffix(pathString, "/") + "/index.html"
	}

	repoURL := fmt.Sprintf("%v/repos/%v/%v/contents/%v", GH_API_BASE, username, repo, pathString)

	resp, err := http.Get(repoURL)
	if err != nil {
		log.Error(err)
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("couldn't read response body: %v", err)
	}

	fileMeta := make(map[string]interface{})
	err = json.Unmarshal(bodyByte, &fileMeta)
	if err != nil {
		log.Errorf("unable to unmarshal json: %v", err)

		// just display content since we can't unmarshal
		w.Header().Add("Content-Type", "application/json")
		w.Write(bodyByte)
		return
	}

	var body []byte
	if content, ok := fileMeta["content"]; ok {
		contentBody, err := base64.StdEncoding.DecodeString(content.(string))
		if err != nil {
			log.Errorf("unable to read content body: %v", err)
		}
		body = contentBody
	} else {

		// just display content since we can't unmarshal
		w.Header().Add("Content-Type", "application/json")
		w.Write(bodyByte)
		return
	}

	extension := filepath.Ext(fileMeta["download_url"].(string))

	mediaType := mime.TypeByExtension(extension)

	w.Header().Add("Content-Type", mediaType)
	w.Write(body)
}
