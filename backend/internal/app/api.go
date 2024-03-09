package app

import (
	"encoding/json"
	"net/http"
)


func (app *Application) RegisterHandlers() {
	http.HandleFunc("/api/images", app.handleListImages)
}

func (app *Application) handleListImages(w http.ResponseWriter, r *http.Request) {
	images, err := app.Kubernetes.ListContainerImages("default") // Assume using "default" namespace
	if err != nil {
		http.Error(w, "Failed to get images", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(images)
}

