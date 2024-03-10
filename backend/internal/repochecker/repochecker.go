package repochecker

import (
	"context"
	"fmt"
	"log"
	"slackwatch/backend/pkg/config"
	"strings"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
)

type Checker struct {
	cfg config.Config // Adjusted to use the top-level Config struct
}

func NewChecker(cfg config.Config) *Checker {
	log.Println("Initializing new Checker with updated configuration")
	return &Checker{cfg: cfg}
}

// getDefaultRepoConfig retrieves the default repository configuration.
func (c *Checker) getDefaultRepoConfig() (config.Repository, error) {
	for _, repo := range c.cfg.Repositories {
		if repo.DefaultRepo {
			log.Printf("Default repository found: %s %v\n", repo.Name, repo.DefaultRepo)
			return repo, nil
		}
	}
	log.Println("No default repository configured")
	return config.Repository{}, fmt.Errorf("no default repository configured")
}

// getRepoConfig retrieves the configuration for a given image name by checking if it matches any configured repository.
// If no specific match is found, it returns the default repository configuration.
func (c *Checker) getRepoConfig(imageName string) (config.Repository, error) {
	for _, repo := range c.cfg.Repositories {
		if strings.Contains(imageName, repo.Name) {
			log.Printf("Repository config found for image: %s, using repository: %s\n", imageName, repo.Name)
			return repo, nil
		}
	}

	// If no specific repository is matched, use the default repository
	return c.getDefaultRepoConfig()
}

func (c *Checker) GetTags(imageName string) ([]string, error) {
	log.Printf("Getting tags for image: %s\n", imageName)

	repoConfig, err := c.getRepoConfig(imageName)
	if err != nil {
		log.Printf("Error getting repository configuration: %v\n", err)
		return nil, err
	}

	var sysCtx types.SystemContext
	if repoConfig.Token != "" && repoConfig.Username != "" && repoConfig.Password != "" {
		// Token, username, and password are provided
		log.Println("Warning: Both token and password are provided; it's recommended to use only one method of authentication.")
		sysCtx = types.SystemContext{
			DockerAuthConfig: &types.DockerAuthConfig{
				IdentityToken: repoConfig.Token,
				Username:      repoConfig.Username,
				Password:      repoConfig.Password,
			},
		}
	} else if repoConfig.Token != "" && repoConfig.Username != "" {
		// Token and username are provided
		sysCtx = types.SystemContext{
			DockerAuthConfig: &types.DockerAuthConfig{
				IdentityToken: repoConfig.Token,
				Username:      repoConfig.Username,
			},
		}
	} else if repoConfig.Token != "" {
		// Only token is provided
		sysCtx = types.SystemContext{
			DockerAuthConfig: &types.DockerAuthConfig{
				IdentityToken: repoConfig.Token,
			},
		}
	} else if repoConfig.Username != "" && repoConfig.Password != "" {
		// Username and password are provided
		sysCtx = types.SystemContext{
			DockerAuthConfig: &types.DockerAuthConfig{
				Username: repoConfig.Username,
				Password: repoConfig.Password,
			},
		}
	} else {
		// No authentication details provided
		log.Println("No authentication details provided, using non-authenticated query")
	}

	// Ensure the imageName includes the registry URL if not already present
	if !strings.Contains(imageName, repoConfig.Name) {
		imageName = repoConfig.Name + "/" + imageName
		log.Printf("Image name updated to: %s\n", imageName)
	}

	ref, err := docker.ParseReference("//" + imageName)
    log.Printf("Parsed reference: %v\n", ref)
	if err != nil {
		log.Printf("Error parsing reference: %v\n", err)
		return nil, fmt.Errorf("error parsing reference: %w", err)
	}
	log.Printf("Fetching tags for image: %s\n", imageName)
	tags, err := docker.GetRepositoryTags(context.Background(), &sysCtx, ref)
	if err != nil {
		log.Printf("Error fetching tags: %v\n", err)
		return nil, fmt.Errorf("error fetching tags: %w", err)
	}

	log.Printf("Tags retrieved for image %s: %v\n", imageName, tags)
	return tags, nil
}
