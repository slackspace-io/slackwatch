package app

import (
	"encoding/json"
	"fmt"
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
    cfg, err := config.LoadConfig("/app/config/config.yaml")
    if err != nil {
        return nil, err
    }
    repoChecker := repochecker.NewChecker(cfg.Repositories)
    k8sClient, err := kubernetes.NewClient(&cfg.Kubernetes, repoChecker)
    if err != nil {
        return nil, err
    }
    
    notificationManager := notifications.NewManager()
    
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

        // First, find containers with the specific annotation
        containers, err := app.Kubernetes.FindContainersWithAnnotation("", "diun.enable", "true")
        fmt.Println(containers)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Then, pass the result directly to CheckForImageUpdates
        updates, err := app.Kubernetes.CheckForImageUpdates(containers)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(updates)
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
