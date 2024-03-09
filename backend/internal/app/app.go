package app

import (
	"encoding/json"
	"log"
	"net/http"
	"slackwatch/backend/internal/kubernetes"
	"slackwatch/backend/internal/notifications"
	"slackwatch/backend/internal/repochecker"
	"slackwatch/backend/pkg/config"
)

type Application struct {
    Kubernetes    *kubernetes.Client
    Notifications *notifications.Manager
    RepoChecker   *repochecker.Checker
}

func Initialize() (*Application, error) {
    // Initialize Kubernetes client, notifications, and repo checker here
    // For example:
    cfg, err := config.LoadConfig("/app/config/config.yaml")
    if err != nil {
        return nil, err
    }
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

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*") // Allow any domain, adjust as necessary for security
}

func (app *Application) setupRoutes() {
    http.HandleFunc("/api/pods", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w) // Enable CORS for this endpoint
        // Assuming you're looking for pods with an annotation "monitoring" set to "enabled"
        pods, err := app.Kubernetes.FindContainersWithAnnotation("", "diun.enable", "true")
        //log result
        log.Println(pods)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(pods)
    })
}

func (app *Application) Run() error {
    log.Println("Starting Slackwatch")
    // Start HTTP server
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
    return nil
}
