package notifications

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"slackwatch/backend/pkg/config" // Import your config package
	"time"
)

type Manager struct {
	ntfyConfig config.NtfyConfig
}

func NewManager(ntfyConfig config.NtfyConfig) *Manager {
	return &Manager{
		ntfyConfig: ntfyConfig,
	}
}

// Example function to send a notification
func (m *Manager) SendNotification(message string) error {
	// Implement the logic to send a notification
	return nil
}

// NotificationPayload defines the structure of the notification payload
type NotificationPayload struct {
	Message  string `json:"message"`
	Priority int    `json:"priority"`
	LastSent string `json:"lastSent"`
}

// SendNtfyContainerUpdate sends a notification to ntfy

// Modify the function signature to return time.Time
func (m *Manager) SendNtfyContainerUpdate(container, currentTag, newTag, foundAt string) (string, error) {
	cfg := m.ntfyConfig
	message := fmt.Sprintf("🔔 *Update Available!* 🔔\n\n*Container:* %s\n*Current Tag:* %s\n*New Tag:* %s\n*Found At:* %s", container, currentTag, newTag, foundAt)
	log.Printf("Sending notification to ntfy: %s", message)
	payload := NotificationPayload{
		Message:  message,
		Priority: cfg.Priority,
	}
	notificationTime, err := m.SendNtfyNotificationGeneric(payload)
	if err != nil {
		return notificationTime, err
	}

	return notificationTime, nil
}

// More Generic ntfy notification
func (m *Manager) SendNtfyNotificationGeneric(payload NotificationPayload) (string, error) {
	cfg := m.ntfyConfig
	log.Printf("Sending notification to ntfy: %s", payload.Message)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", cfg.URL, cfg.Topic), bytes.NewBuffer([]byte(payload.Message)))
	if err != nil {
		return time.Now().Format(time.RFC3339), err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return time.Now().Format(time.RFC3339), err
	}
	log.Printf("Sent notification to ntfy, status code: %d", resp.StatusCode)
	defer resp.Body.Close()
	// Get the current time in parsable format
	notificationTime := time.Now().Format(time.RFC3339)

	if resp.StatusCode != http.StatusOK {
		return notificationTime, fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
	}

	// Get the current time
	//ensure time is in time.Time parseable format

	return notificationTime, nil
}
