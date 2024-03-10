package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slackwatch/backend/pkg/config" // Import your config package
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
}

// SendNtfyNotification sends a notification to ntfy
func (m *Manager) SendNtfyNotification(container, currentTag, newTag, foundAt string) error {
    cfg := m.ntfyConfig // Now correctly accessing the configuration
    payload := NotificationPayload{
        Message:  fmt.Sprintf("Container: %s, Current Tag: %s, New Tag: %s, Found At: %s", container, currentTag, newTag, foundAt),
        Priority: cfg.Priority,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    log.Printf("Sending notification to ntfy: %s", payloadBytes)
    req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", cfg.URL, cfg.Topic), bytes.NewBuffer(payloadBytes))
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer "+cfg.Token)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    log.Printf("Sent notification to ntfy, status code: %d", resp)
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
    }
    return nil
}

