package app

import (
	"context"
	"log"
	"net/http"
	"slackwatch/backend/internal/kubernetes"
	"slackwatch/backend/internal/notifications"
	"slackwatch/backend/internal/repochecker"
	"slackwatch/backend/internal/service"
	"slackwatch/backend/internal/storage"
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

    // Initialize services and storage
    fileStorage := &storage.FileStorage{FilePath: "data.json"}
    dataService := &service.DataService{Storage: fileStorage}

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

    return app, nil
}

func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (app *Application) setupRoutes() {
    // Setup routes for data service with CORS enabled
    http.HandleFunc("/api/images", app.enableCorsMiddleware(app.DataService.HandleGetData))
    http.HandleFunc("/api/refresh", app.enableCorsMiddleware(app.handleRefreshData))

    // Setup other routes as before
    // Example: http.HandleFunc("/api/pods", app.handlePods)
}

// enableCorsMiddleware wraps an http.HandlerFunc with CORS headers
func (app *Application) enableCorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        enableCors(w)
        next.ServeHTTP(w, r)
    }
}

func (app *Application) handleRefreshData(w http.ResponseWriter, r *http.Request) {
    // Logic to refresh data manually
    log.Println("Manually refreshing data...")
    app.runScheduledTask() // Directly invoke the task logic
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
    updates, err := app.Kubernetes.CheckForImageUpdates(containers)
    if err != nil {
        log.Printf("Error checking for image updates: %v", err)
        return
    }
    // Save updates data locally
    log.Printf("Saving updates data: %v", updates)
    err = app.DataService.SaveData(context.Background(), updates)
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

// Include other methods (e.g., scheduleKubernetesTasks) as they were previously defined
