package app

import (
	"encoding/json"
	"net/http"
	"slackwatch/backend/internal/kubernetes"
	"slackwatch/backend/internal/notifications"
	"slackwatch/backend/internal/repochecker"
	"slackwatch/backend/pkg/config"
	"time"
)

type Application struct {
    Kubernetes    *kubernetes.Client
    Notifications *notifications.Manager
    RepoChecker   *repochecker.Checker
}

func Initialize() (*Application, error) {
    // Initialize Kubernetes client, notifications, and repo checker here
    // For example:
    cfg, err := config.LoadConfig("config/config.yaml")
    k8sClient, err := kubernetes.NewClient(&cfg.Kubernetes)
    if err != nil {
        return nil, err
    }
    
    notificationManager := notifications.NewManager()
    repoChecker := repochecker.NewChecker()
    
    app := &Application{
        Kubernetes:    k8sClient,
        Notifications: notificationManager,
        RepoChecker:   repoChecker,
    }
    
    app.setupRoutes()
    
    return app, nil
}

func (app *Application) setupRoutes() {
    http.HandleFunc("/api/pods", func(w http.ResponseWriter, r *http.Request) {
        images, err := app.Kubernetes.ListContainerImages("default")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(images)
    })
}

func (app *Application) Run() error {
    // Your application's main logic
    // For example, monitor Kubernetes pods, check for updates, and send notifications
    
    return nil
}
