package main

import (
    "fmt"
    "slackwatch/backend/internal/kubernetes"
    "slackwatch/backend/pkg/config"
)

func main() {
    cfg, err := config.LoadConfig("config/config.yaml")
    
    if err != nil {
        fmt.Println("Error loading config:", err)
        return
    }

    k8sClient, err := kubernetes.NewClient(&cfg.Kubernetes)
    
    if err != nil {
        fmt.Println("Error creating Kubernetes client:", err)
        return
    }

    // Use k8sClient here
    images, err := k8sClient.ListContainerImages("default") // Assuming you want to list containers in the "default" namespace
    if err != nil {
        fmt.Println("Error listing container images:", err)
        return
    }
    fmt.Println("Running container images:", images)
}
