package kubernetes

import (
    "context"
    "fmt"
    "path/filepath"
    "slackwatch/backend/pkg/config" // Import your config package

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)


type Client struct {
    clientSet *kubernetes.Clientset
}

func NewClient(cfg *config.KubernetesConfig) (*Client, error) {
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

    return &Client{clientSet: clientSet}, nil
}


// ListContainerImages lists all container images in a given namespace
func (c *Client) ListContainerImages(namespace string) ([]string, error) {
    podList, err := c.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
    }

    var images []string
    for _, pod := range podList.Items {
        for _, container := range pod.Spec.Containers {
            images = append(images, container.Image)
        }
    }

    return images, nil
}
