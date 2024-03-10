package app

import (
	"net/http"
)

type SystemConfig struct {
	Schedule string `yaml:"schedule"`
}

func (app *Application) RegisterHandlers() {
	http.HandleFunc("/api/images", app.handleListImages)
	http.HandleFunc("/api/containers", app.enableCorsMiddleware(app.handleContainerInfo))
	http.HandleFunc("/api/imageUpdates", app.enableCorsMiddleware(app.handleImageUpdates))
}
