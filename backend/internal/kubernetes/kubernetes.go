package kubernetes

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"slackwatch/backend/internal/repochecker" // Assuming the import path for repochecker
	"slackwatch/backend/pkg/config"           // Import your config package

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
func (c *Client) FindStatefulSets(namespace string, name string) (map[string]string, error) {
	statefulSets, err := c.clientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list stateful sets in namespace %s: %w", namespace, err)
	}
	// Print each stateful set and metadata
	statefulSetMap := make(map[string]string)
	statefulSetMap["name"] = statefulSets.Name
	statefulSetMap["namespace"] = statefulSets.Namespace
	statefulSetMap["timeScanned"] = time.Now().Format(time.RFC3339)
	//get imagename
	containers := statefulSets.Spec.Template.Spec.Containers
	for _, container := range containers {
		statefulSetMap["image"] = container.Image
	}
	//log all values in map
	log.Printf("Found stateful set: %v", statefulSetMap)
	return statefulSetMap, nil
}

// Get single container info
func (c *Client) GetContainerInfo(namespace string, podName string, containerName string) (map[string]string, error) {
	pod, err := c.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s in namespace %s: %w", podName, namespace, err)
	}

	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			return map[string]string{
				"podName":       pod.Name,
				"containerName": container.Name,
				"image":         container.Image,
				"timeScanned":   time.Now().Format(time.RFC3339),
			}, nil
		}
	}

	return nil, fmt.Errorf("container %s not found in pod %s", containerName, podName)

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
				excludePattern, _ := pod.ObjectMeta.Annotations["slackwatch.exclude"]
				includePattern, _ := pod.ObjectMeta.Annotations["slackwatch.include"]
				gitopsRepo, _ := pod.ObjectMeta.Annotations["slackwatch.repo"]
				containersWithAnnotation = append(containersWithAnnotation, map[string]string{
					"podName":        pod.Name,
					"containerName":  container.Name,
					"image":          container.Image,
					"timeScanned":    time.Now().Format(time.RFC3339),
					"excludePattern": excludePattern,
					"includePattern": includePattern,
					"gitopsRepo":     gitopsRepo,
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
			isNewer, err := c.isTagNewer(currentTag, newTag, container["excludePattern"], container["includePattern"])
			if err != nil {
				log.Printf("Error comparing versions: %v", err)
				continue
			}
			if isNewer {
				log.Printf("New tag found %s", newTag)
				container["currentTag"] = currentTag                   // Add current tag to the container's map
				container["newTag"] = newTag                           // Add new tag information
				container["isNewer"] = strconv.FormatBool(isNewer)     // Add isNewer flag
				container["foundAt"] = time.Now().Format(time.RFC3339) // Add foundAt timestamp
				updates = append(updates, container)                   // Append the modified container map
			}
		}
	}
	return updates, nil
}

// isTagNewer compares two semantic versioning tags and checks if the newTag is actually newer than the currentTag.
// It also checks if the newTag matches any of the exclude patterns provided in the config or the specific container exclude pattern.
func (c *Client) isTagNewer(currentTag, newTag, excludePattern, includePattern string) (bool, error) {
	if includePattern != "" {
		// Convert wildcard pattern to valid regex pattern
		regexPattern := strings.ReplaceAll(includePattern, "*", ".*")
		matched, err := regexp.MatchString(regexPattern, newTag)
		if err != nil {
			return false, fmt.Errorf("error matching include pattern: %w", err)
		}
		if !matched {
			log.Printf("Skipping tag %s due to not matching include pattern %s", newTag, includePattern)
			return false, nil // If newTag does not match the include pattern, it's not considered.
		}
	}

	includePatterns := c.config.Magic.IncludePatterns
	if len(includePatterns) > 0 {
		includeMatched := false
		for _, pattern := range includePatterns {
			regexPattern := strings.ReplaceAll(pattern, "*", ".*")
			matched, err := regexp.MatchString(regexPattern, newTag)
			if err != nil {
				return false, fmt.Errorf("error matching include pattern: %w", err)
			}
			if matched {
				includeMatched = true
				break
			}
		}
		if !includeMatched {
			log.Printf("Tag %s does not match any include pattern", newTag)
			return false, nil
		}
	}

	if excludePattern != "" {
		// Convert wildcard pattern to valid regex pattern
		regexPattern := strings.ReplaceAll(excludePattern, "*", ".*")
		matched, err := regexp.MatchString(regexPattern, newTag)
		if err != nil {
			return false, fmt.Errorf("error matching exclude pattern: %w", err)
		}
		if matched {
			log.Printf("Skipping tag %s due to exclude pattern %s", newTag, excludePattern)
			return false, nil // If newTag matches the exclude pattern, it's not considered newer.
		}
	}

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
