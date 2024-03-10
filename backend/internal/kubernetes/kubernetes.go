package kubernetes

import (
    "context"
    "fmt"
    "path/filepath"
    "slackwatch/backend/pkg/config" // Import your config package
    "time"
    "strings"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    corev1 "k8s.io/api/core/v1"
    "slackwatch/backend/internal/repochecker" // Assuming the import path for repochecker
)


type Client struct {
    clientSet *kubernetes.Clientset
    repoChecker *repochecker.Checker
}

func NewClient(cfg *config.KubernetesConfig, checker *repochecker.Checker) (*Client, error) {
    var k8sConfig *rest.Config
    var err error

    if cfg.UseInClusterConfig {
        k8sConfig, err = rest.InClusterConfig()
        if err != nil {
            return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
        }
    } else {
        kubeconfigPath := filepath.Clean(cfg.OutOfClusterConfig.KubeconfigPath)
        k8sConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
        if err != nil {
            return nil, fmt.Errorf("failed to get kubeconfig from path %s: %w", kubeconfigPath, err)
        }
    }

    clientSet, err := kubernetes.NewForConfig(k8sConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create clientset: %w", err)
    }

    return &Client{clientSet: clientSet, repoChecker: checker}, nil
}

// FindContainersWithAnnotation finds all containers in a given namespace (or all namespaces if namespace is empty) that have a specific metadata annotation
func (c *Client) FindContainersWithAnnotation(namespace string, annotationKey string, annotationValue string) ([]map[string]string, error) {
    // If namespace is provided, search within that namespace. Otherwise, search across all namespaces.
    podList, err := c.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
    }

    var containersWithAnnotation []map[string]string
    for _, pod := range podList.Items {
        for _, container := range pod.Spec.Containers {
            if value, ok := pod.ObjectMeta.Annotations[annotationKey]; ok && value == annotationValue {
                containersWithAnnotation = append(containersWithAnnotation, map[string]string{
                    "podName":       pod.Name,
                    "containerName": container.Name,
                    "image":         container.Image,
                    "timeScanned":   time.Now().Format(time.RFC3339),
                })
            }
        }
    }

    return containersWithAnnotation, nil
}

// ListContainerImages lists all container images in a given namespace or all namespaces if namespace is empty
func (c *Client) ListContainerImages(namespace string) ([]map[string]string, error) {
    var podList *corev1.PodList
    var err error

    if namespace == "" {
        podList, err = c.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
    } else {
        podList, err = c.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
    }

    if err != nil {
        return nil, fmt.Errorf("failed to list pods in namespace '%s': %w", namespace, err)
    }

    var images []map[string]string
    for _, pod := range podList.Items {
        for _, container := range pod.Spec.Containers {
            images = append(images, map[string]string{
                "image": container.Image,
                "timeScanned": time.Now().Format(time.RFC3339),
            })
        }
    }

    return images, nil
}
// CheckForImageUpdates checks for newer tags of the images
func (c *Client) CheckForImageUpdates(containers []map[string]string) ([]map[string]string, error) {
    var updates []map[string]string
    for _, container := range containers {
        fmt.Println(container)
        imageName := container["image"]
        fmt.Println(imageName)
        // Assuming the image name is in the format "repo/image:tag"
        parts := strings.Split(imageName, ":")
        if len(parts) != 2 {
            continue // Skip if the format is unexpected
        }
        repo, currentTag := parts[0], parts[1]
        fmt.Println(repo)
        fmt.Println(currentTag)
        newTags, err := c.repoChecker.GetTags(repo)
        if err != nil {
            continue // Skip on error
        }
        for _, newTag := range newTags {
            if newTag != "" && newTag != currentTag {
                updates = append(updates, map[string]string{
                    "image": imageName,
                    "currentTag": currentTag,
                    "newTag": newTag,
                })
            }
        }
    }
    fmt.Println(updates)
    return updates, nil
}
