package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "file error", 400)
		return
	}
	defer file.Close()

	eventId := r.FormValue("eventId")
	if eventId == "" {
		http.Error(w, "missing eventId", 400)
		return
	}

	userId := r.FormValue("userId") // <-- Get the userId from the form
	timestamp := time.Now().Unix()

	filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
	dir := filepath.Join("uploads", eventId)
	os.MkdirAll(dir, 0755)

	// Save the image file
	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		http.Error(w, "save error", 500)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	// Save metadata as JSON next to image
	meta := map[string]string{
		"userId":    userId,
		"timestamp": fmt.Sprintf("%d", timestamp),
	}

	metaPath := filepath.Join(dir, filename+".json")
	metaFile, err := os.Create(metaPath)
	if err == nil {
		defer metaFile.Close()
		json.NewEncoder(metaFile).Encode(meta)
	}

	w.Write([]byte("ok"))
}

func HandleImages(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		http.Error(w, "missing eventId", http.StatusBadRequest)
		return
	}

	dir := filepath.Join("uploads", eventId)
	files, err := os.ReadDir(dir)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{})
		return
	}

	// Sort by timestamp prefix in filename (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})

	type ImageEntry struct {
		URL    string `json:"url"`
		UserId string `json:"userId"`
		Time   string `json:"timestamp"`
	}

	var entries []ImageEntry

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) == ".json" {
			continue
		}

		filename := file.Name()
		metaPath := filepath.Join(dir, filename+".json")

		var userID string
		var timestamp string

		// Try reading metadata
		if fmeta, err := os.Open(metaPath); err == nil {
			defer fmeta.Close()
			var meta map[string]string
			if err := json.NewDecoder(fmeta).Decode(&meta); err == nil {
				userID = meta["userId"]
				timestamp = meta["timestamp"]
			}
		}

		url := fmt.Sprintf("/uploads/%s/%s", eventId, filename)
		entries = append(entries, ImageEntry{
			URL:    url,
			UserId: userID,
			Time:   timestamp,
		})

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func HandleEventConfig(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		http.Error(w, "missing eventId", 400)
		return
	}

	filePath := filepath.Join("events", eventId+".json")
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "event not found", 404)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, f)
}
