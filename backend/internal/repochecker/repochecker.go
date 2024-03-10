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
	cfg []config.Repository // Configuration for multiple repositories
}

func NewChecker(cfg []config.Repository) *Checker {
	log.Println("Initializing new Checker")
	return &Checker{cfg: cfg}
}

// getRepoConfig retrieves the configuration for a given registry URL.
func (c *Checker) getRepoConfig(registryURL string) (config.Repository, error) {
	log.Printf("Retrieving repo config for registry URL: %s\n", registryURL)
	for _, repo := range c.cfg {
		if strings.Contains(registryURL, repo.Name) {
			log.Printf("Config found for registry: %s\n", registryURL)
			return repo, nil
		}
	}
	log.Printf("No config found for registry: %s\n", registryURL)
	return config.Repository{}, fmt.Errorf("no config found for registry: %s", registryURL)
}

func (c *Checker) GetTags(imageName string) ([]string, error) {
	log.Printf("Getting tags for image: %s\n", imageName)
	// Assuming GHCR is part of the imageName, adjust if necessary
	repoConfig, err := c.getRepoConfig("ghcr.io")
	if err != nil {
		log.Printf("Error getting repo config: %v\n", err)
		return nil, err
	}

	var sysCtx types.SystemContext
	if repoConfig.Token != "" {
		sysCtx = types.SystemContext{
			DockerAuthConfig: &types.DockerAuthConfig{
				IdentityToken: repoConfig.Token,
			},
		}
	} else {
		sysCtx = types.SystemContext{
			DockerAuthConfig: &types.DockerAuthConfig{
				Username: repoConfig.Username,
				Password: repoConfig.Password,
			},
		}
	}

	ref, err := docker.ParseReference("//" + imageName)
	if err != nil {
		log.Printf("Error parsing reference: %v\n", err)
		return nil, fmt.Errorf("error parsing reference: %w", err)
	}
	tags, err := docker.GetRepositoryTags(context.Background(), &sysCtx, ref)
	if err != nil {
		log.Printf("Error fetching tags: %v\n", err)
		return nil, fmt.Errorf("error fetching tags: %w", err)
	}

	log.Printf("Tags retrieved for image %s: %v\n", imageName, tags)
	return tags, nil
}
