package app

import (
	"context"
	"encoding/json"
	"net/http"
)

type SystemConfig struct {
	Schedule string `yaml:"schedule"`
}

func (app *Application) RegisterHandlers() {
	http.HandleFunc("/api/images", app.handleListImages)
}

func (app *Application) handleListImages(w http.ResponseWriter, r *http.Request) {
	refresh := r.URL.Query().Get("refresh")
	if refresh == "true" {
		// Trigger data refresh logic here
	} else {
		// Use saved data to respond to the request
		savedData, err := app.DataService.GetData(context.Background())
		if err != nil {
			http.Error(w, "Failed to get saved data", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(savedData)
	}
}

