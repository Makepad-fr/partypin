package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Makepad-fr/photo-pwa/internal"
)

type EventConfig struct {
	Title             string `json:"title"`
	AllowGallery      bool   `json:"allowGallery"`
	AllowAny          bool   `json:"allowAny"`
	RequireTakenToday bool   `json:"requireTakenToday"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/create-event", handleCreateEvent)
	http.HandleFunc("/event-config", internal.HandleEventConfig)

	http.HandleFunc("/upload", internal.HandleUpload)
	http.HandleFunc("/images", internal.HandleImages)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	fmt.Println("Listening on http://localhost:6060")
	http.ListenAndServe("0.0.0.0:6060", nil)
}
func handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config EventConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	eventId := internal.GenerateEventID(6)
	dir := "events"
	os.MkdirAll(dir, 0755)

	filePath := filepath.Join(dir, eventId+".json")
	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "could not save config", 500)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(config); err != nil {
		http.Error(w, "failed to write file", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"eventId": eventId,
	})
}
