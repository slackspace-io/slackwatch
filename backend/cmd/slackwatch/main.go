package main

import (
    "fmt"
    "slackwatch/backend/internal/kubernetes"
    "slackwatch/backend/pkg/config"
    "time"
)

func main() {
    cfg, err := config.LoadConfig("config/config.yaml")
    
    if err != nil {
        fmt.Println("Error loading config:", err)
        return
    }

    if cfg.Kubernetes.PollingInterval <= 0 {
        fmt.Println("PollingInterval must be a positive integer")
        return
    }

    k8sClient, err := kubernetes.NewClient(&cfg.Kubernetes)
    
    if err != nil {
        fmt.Println("Error creating Kubernetes client:", err)
        return
    }

    ticker := time.NewTicker(time.Duration(cfg.Kubernetes.PollingInterval) * time.Second)
    defer ticker.Stop()

    go func() {
        for {
            select {
            case <-ticker.C:
                images, err := k8sClient.ListContainerImages("default")
                if err != nil {
                    fmt.Println("Error listing container images:", err)
                    continue
                }
                fmt.Println("Running container images:", images)
            }
        }
    }()

    // Block main goroutine indefinitely
    select {}
}
