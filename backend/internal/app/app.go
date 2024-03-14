package app

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"slackwatch/backend/internal/kubernetes"
	"slackwatch/backend/internal/notifications"
	"slackwatch/backend/internal/repochecker"
	"slackwatch/backend/internal/service"
	"slackwatch/backend/pkg/config"
	"strconv"
	"time"
)

type Application struct {
	Kubernetes    *kubernetes.Client
	Notifications *notifications.Manager
	RepoChecker   *repochecker.Checker
	Scheduler     *cron.Cron
	DataService   *service.DataService
	System        config.SystemConfig `yaml:"system"`
	Config        *config.Config      // Add this line to hold the entire configuration
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
	k8sClient, err := kubernetes.NewClient(&cfg.Kubernetes, repoChecker, cfg)
	if err != nil {
		log.Printf("Failed to initialize Kubernetes client: %v", err)
		return nil, err
	}

	// Modify this line to pass the configuration to NewManager
	notificationManager := notifications.NewManager(cfg.Notifications.Ntfy)

	app := &Application{
		Kubernetes:    k8sClient,
		Notifications: notificationManager,
		RepoChecker:   repoChecker,
		DataService:   dataService,
		System:        cfg.System,
		Config:        cfg, // Store the entire configuration in the Application struct
	}

	app.setupRoutes()
	app.scheduleTasks()
	if app.System.RunAtStartup { // Check if runAtStartup is true
		go app.runScheduledTask()
	}
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
	http.HandleFunc("/api/data/combined", app.enableCorsMiddleware(app.handleCombinedData))
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

// create test function to run at startup
func (app *Application) getCombinedData() {
	log.Println("Running at startup...")
	combinedData, err := app.DataService.GetCombinedData(context.Background())
	//convert to json
	if err != nil {
		log.Printf("Error getting previous data: %v", err)
	}
	//log the previous data
	log.Printf("Previous data: %v", combinedData)
}

// runScheduledTask encapsulates the logic to be executed on schedule.
func (app *Application) runScheduledTask() {
	log.Println("Scheduled task running...")
	//get previous data
	combinedData, err := app.DataService.GetCombinedData(context.Background())
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
	for _, update := range updates {
		var lastNotificationTime = time.Time{}
		for _, data := range combinedData {
			if data["containerName"] == update["containerName"] {
				lastSentTime, ok := data["sentTime"].(string)
				if !ok {
					continue
				}
				lastNotificationTime, err = time.Parse(time.RFC3339, lastSentTime)
				if err != nil {
					log.Printf("Error parsing time: %v", err)
					continue
				}
				break
			}
		}

		//Compare the last notification time with the current time and the reminder interval
		reminderDuration, err := convertReminderToDuration(app.Config.Notifications.Ntfy.Reminder)
		if err != nil {
			log.Printf("Error converting reminder to duration: %v", err)
			continue
		}
		if time.Since(lastNotificationTime) >= reminderDuration {
			notificationTime, err := app.Notifications.SendNtfyNotification(update["containerName"], update["currentTag"], update["newTag"], update["foundAt"])
			//add sentTime to updates object
			update["sentTime"] = notificationTime
			if err != nil {
				log.Printf("Failed to send notification for container %s: %v", update["containerName"], err)
			}
		} else {
			log.Printf("Skipping notification for container %s: reminder interval not reached", update["containerName"])
		}
	}
	err = app.DataService.SaveData(context.Background(), "image", updates)
	if err != nil {
		log.Printf("Data saved successfully")
		log.Printf("Error saving image data: %v", err)
	}
}

func convertReminderToDuration(reminder string) (time.Duration, error) {
	// Get the unit of time and the value
	unit := reminder[len(reminder)-1:]
	value, err := strconv.Atoi(reminder[:len(reminder)-1])
	if err != nil {
		return 0, fmt.Errorf("invalid reminder format: %v", err)
	}

	// Convert the value to a duration based on the unit of time
	switch unit {
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "w":
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown unit of time: %s", unit)
	}
}

func (app *Application) Run() error {
	log.Println("Starting the application on port :8080")
	// Start HTTP server in a goroutine so that it doesn't block
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Send a fake notification with spoof data
	notificationTime, err := app.Notifications.SendNtfyNotification("spoof-container", "v1.0.0", "v1.0.1", "2023-04-01T12:00:00Z")
	log.Printf("Fake notification sent at: %s", notificationTime)
	if err != nil {
		log.Printf("Failed to send fake notification: %v", err)
	} else {
		log.Println("Fake notification sent successfully")
	}
	// Keep the application running
	select {}
}
