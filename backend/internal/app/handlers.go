package app

import (
	"context"
	"encoding/json"
	"net/http"
)

// handleContainerInfo responds with the container data.
func (app *Application) handleContainerInfo(w http.ResponseWriter, r *http.Request) {
	containerData, err := app.DataService.GetData(context.Background(), "containers.json")
	if err != nil {
		http.Error(w, "Failed to get container data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(containerData)
}

// handleImageUpdates responds with the image update data.
func (app *Application) handleImageUpdates(w http.ResponseWriter, r *http.Request) {
	imageData, err := app.DataService.GetData(context.Background(), "imageUpdates.json")
	if err != nil {
		http.Error(w, "Failed to get image data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(imageData)
}

// handleListImages handles requests for listing images.
func (app *Application) handleListImages(w http.ResponseWriter, r *http.Request) {
	refresh := r.URL.Query().Get("refresh")
	if refresh == "true" {
		// Trigger data refresh logic here
	} else {
		// Use saved data to respond to the request
		savedData, err := app.DataService.GetData(context.Background(), "images.json")
		if err != nil {
			http.Error(w, "Failed to get saved data", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(savedData)
	}
}