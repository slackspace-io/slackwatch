package kubernetes

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"slackwatch/backend/pkg/config" // Import your config package
	"slackwatch/backend/internal/repochecker" // Assuming the import path for repochecker

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientSet   *kubernetes.Clientset
	repoChecker *repochecker.Checker
	config      *config.Config
}

func NewClient(cfg *config.KubernetesConfig, checker *repochecker.Checker, appConfig *config.Config) (*Client, error) {
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

	return &Client{clientSet: clientSet, repoChecker: checker, config: appConfig}, nil
}

// FindContainersWithAnnotation finds all containers in a given namespace (or all namespaces if namespace is empty) that have a specific metadata annotation
func (c *Client) FindContainersWithAnnotation(namespace string, annotationKey string, annotationValue string) ([]map[string]string, error) {
	// If namespace is provided, search within that namespace. Otherwise, search across all namespaces.
	podList, err := c.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	log.Printf("Found %d pods", len(podList.Items))
	// Print each pod and metadata

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
	log.Printf("Found %d containers with annotation", len(containersWithAnnotation))
	for _, container := range containersWithAnnotation {
		log.Printf("Container: %v", container)
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
				"image":       container.Image,
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
		imageName := container["image"]
		// Assuming the image name is in the format "repo/image:tag"
		parts := strings.Split(imageName, ":")
		if len(parts) != 2 {
			continue // Skip if the format is unexpected
		}
		repo, currentTag := parts[0], parts[1]
		newTags, err := c.repoChecker.GetTags(repo)
		if err != nil {
			continue // Skip on error
		}
		for _, newTag := range newTags {
			isNewer, err := c.isTagNewer(currentTag, newTag)
			if err != nil {
				log.Printf("Error comparing versions: %v", err)
				continue
			}
			if isNewer {
				log.Printf("New tag found %s", newTag)
				container["currentTag"] = currentTag // Add current tag to the container's map
				container["newTag"] = newTag         // Add new tag information
				container["isNewer"] = strconv.FormatBool(isNewer) // Add isNewer flag
				container["foundAt"] = time.Now().Format(time.RFC3339) // Add foundAt timestamp
				updates = append(updates, container) // Append the modified container map
			}
		}
	}
	return updates, nil
}

// isTagNewer compares two semantic versioning tags and checks if the newTag is actually newer than the currentTag.
// It also checks if the newTag matches any of the exclude patterns provided in the config.
func (c *Client) isTagNewer(currentTag, newTag string) (bool, error) {
    excludePatterns := c.config.Magic.ExcludePatterns
    log.Printf("Exclude patterns: %v", excludePatterns)
    if len(excludePatterns) > 0 {
        for _, pattern := range excludePatterns {
            // Convert wildcard pattern to valid regex pattern
            regexPattern := strings.ReplaceAll(pattern, "*", ".*")
            matched, err := regexp.MatchString(regexPattern, newTag)
            if err != nil {
                return false, fmt.Errorf("error matching exclude pattern: %w", err)
            }
            if matched {
                log.Printf("Skipping tag %s due to exclude pattern %s", newTag, pattern)
                return false, nil // If newTag matches any exclude pattern, it's not considered newer.
            }
        }
    }

    
    // Parse semantic versions
    currentVersion, err := semver.NewVersion(currentTag)
    if err != nil {
        return false, fmt.Errorf("error parsing current tag '%s' as semver: %w", currentTag, err)
    }
    newVersion, err := semver.NewVersion(newTag)
    if err != nil {
        return false, fmt.Errorf("error parsing new tag '%s' as semver: %w consider updating your exclude patterns", newTag, err)
    }

    // Compare versions
    return newVersion.GreaterThan(currentVersion), nil
}