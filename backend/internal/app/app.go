package app

import (
	"context"
	"log"
	"net/http"
	"slackwatch/backend/internal/kubernetes"
	"slackwatch/backend/internal/notifications"
	"slackwatch/backend/internal/repochecker"
	"slackwatch/backend/internal/service"
	"slackwatch/backend/pkg/config"

	"github.com/robfig/cron/v3"
)

type Application struct {
    Kubernetes    *kubernetes.Client
    Notifications *notifications.Manager
    RepoChecker   *repochecker.Checker
    Scheduler     *cron.Cron
    DataService   *service.DataService
    System        config.SystemConfig `yaml:"system"`
}

func Initialize() (*Application, error) {
    cfg, err := config.LoadConfig("/app/config/config.yaml")
    if err != nil {
        log.Printf("Failed to load configuration: %v", err)
        return nil, err
    }

    // Initialize services
    dataService := &service.DataService{}

    // Initialize other components as before
    repoChecker := repochecker.NewChecker(*cfg)
    k8sClient, err := kubernetes.NewClient(&cfg.Kubernetes, repoChecker)
    if err != nil {
        log.Printf("Failed to initialize Kubernetes client: %v", err)
        return nil, err
    }
    notificationManager := notifications.NewManager()

    app := &Application{
        Kubernetes:    k8sClient,
        Notifications: notificationManager,
        RepoChecker:   repoChecker,
        DataService:   dataService,
        System:        cfg.System,
    }

    app.setupRoutes()
    app.scheduleTasks()
    go app.runScheduledTask()
    return app, nil
}

func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (app *Application) setupRoutes() {
    http.HandleFunc("/api/containers", app.enableCorsMiddleware(app.handleContainerInfo))
    http.HandleFunc("/api/imageUpdates", app.enableCorsMiddleware(app.handleImageUpdates))
    // Add the new route for data refresh
    http.HandleFunc("/api/data/refresh", app.enableCorsMiddleware(app.handleDataRefresh))
}

// handleDataRefresh triggers the scheduled task manually.
func (app *Application) handleDataRefresh(w http.ResponseWriter, r *http.Request) {

    go app.runScheduledTask() // Run in a goroutine to not block the HTTP response

    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("Scheduled task triggered"))
}

// enableCorsMiddleware wraps an http.HandlerFunc with CORS headers
func (app *Application) enableCorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        if r.Method == "OPTIONS" {
            return
        }
        next.ServeHTTP(w, r)
    }
}

func (app *Application) scheduleTasks() {
    app.Scheduler = cron.New()
    _, err := app.Scheduler.AddFunc(app.System.Schedule, app.runScheduledTask)
    if err != nil {
        log.Fatalf("Failed to schedule tasks: %v", err)
    }
    app.Scheduler.Start()
}

// runScheduledTask encapsulates the logic to be executed on schedule.
func (app *Application) runScheduledTask() {
    log.Println("Scheduled task running...")
    containers, err := app.Kubernetes.FindContainersWithAnnotation("", "slackwatch.enable", "true")
    if err != nil {
        log.Printf("Error finding containers: %v", err)
        return
    }
    log.Printf("Found %d containers", len(containers))
    for _, container := range containers {
        log.Printf("Container: %v", container)
    }
    err = app.DataService.SaveData(context.Background(), "container", containers)
    if err != nil {
        log.Printf("Error saving container data: %v", err)
    }

    updates, err := app.Kubernetes.CheckForImageUpdates(containers)
    if err != nil {
        log.Printf("Error checking for image updates: %v", err)
        return
    }
    err = app.DataService.SaveData(context.Background(), "image", updates)
    if err != nil {
        log.Printf("Error saving updates data: %v", err)
    }
    log.Printf("Data saved successfully")
}

func (app *Application) Run() error {
    log.Println("Starting the application on port :8080")
    // Start HTTP server
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
    return nil
}
