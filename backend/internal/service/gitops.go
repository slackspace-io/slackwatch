package service

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"slackwatch/backend/internal/notifications"
	"slackwatch/backend/pkg/config"
	"strings"
	"time"
)

type Gitops struct {
	cfg config.Config
	// Add a field to hold the entire configuration
}

func NewGitOps(cfg config.Config) *Gitops {
	log.Println("Initializing new GitOps with updated configuration")
	return &Gitops{cfg: cfg}
}

func (gs *Gitops) UpdateRequest(updateRequest map[string]string, runningData map[string]string) {
	log.Printf("Received update: %v", updateRequest)
	log.Printf("Running data: %v", runningData)
	// Check if the new tag is the same as the current tag
	if gs.SanityChecks(updateRequest, runningData) {
		// Clone the repository
		log.Printf("Cloning repository %s", updateRequest["repo"])
	}
	repoConfig, err := gs.GetRepoConfig(updateRequest["repo"])
	if err != nil {
		log.Printf("Error getting repository configuration: %v", err)
		return
	}
	//log config available
	log.Printf("Repository config: %v", repoConfig)
	// Clone the repository
	// Delete first
	gs.DeleteLocalRepo()
	log.Printf("Cloning repository %s", repoConfig.RepoURL)
	repo, err := gs.CloneRepo(repoConfig)
	if err != nil {
		gs.DeleteLocalRepo()
		log.Printf("Error cloning repository: %v", err)
		return
	}
	log.Printf("Repository cloned successfully: %s", repoConfig.RepoURL)
	branches, err := repo.Branches()
	if err != nil {
		log.Printf("Error getting branches: %v", err)
		return
	}
	log.Printf("Branches: %v", branches)
	//list metadata of repo cloned
	//list files
	//cleanup
	fullPath := gs.BuildPath(repo, updateRequest)
	//send notification
	gs.UpdateYamls(repo, updateRequest, fullPath)
	gs.CheckForChangesAndCommit(repoConfig, repo, updateRequest, runningData)

}

func (gs *Gitops) UpdateYamls(repo *git.Repository, updateRequest map[string]string, fullPath string) {
	newTag := updateRequest["newTag"]
	currentTag := updateRequest["currentTag"]
	currentImage := updateRequest["image"]
	log.Printf("Updating YAML files for %s from %s to %s ", updateRequest["repo"], currentTag, newTag)
	// Update the YAML files
	//log config available
	log.Printf("Updating YAML files for %s", updateRequest["repo"])
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking path: %v", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}
		log.Printf("Updating file: %s", path)
		gs.ProcessYaml(repo, path, currentImage, newTag, currentTag)
		return nil
		//list files
	})
	if err != nil {
		log.Printf("Error updating YAML files: %v", err)
	}
}

func (gs *Gitops) ProcessYaml(repo *git.Repository, path, currentImage, newTag string, currentTag string) {
	// Process the YAML file
	//log config available
	log.Printf("Processing YAML file: %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}
	var obj map[string]interface{}
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		log.Printf("Error unmarshalling YAML: %v", err)
		return
	}
	//print fields from obj
	//get image
	//check if statefulset or deployment
	if kind, ok := obj["kind"]; ok {
		if kind == "StatefulSet" || kind == "Deployment" {
			// Update the StatefulSet
			log.Printf("Found StatefulSet or Deployment: %s in %s", kind, path)
			obj, err := gs.FindAndUpdateYaml(obj, currentImage, currentTag, newTag)
			if err != nil {
				log.Printf("Not Updating YAML: %v", err)
				return
			}
			modifiedYaml, err := yaml.Marshal(obj)
			if err != nil {
				log.Printf("Error marshalling YAML: %v", err)
				return
			}
			err = os.WriteFile(path, modifiedYaml, 0644)
			if err != nil {
				log.Printf("Error writing to file: %v", err)
				return
			}
			log.Printf("YAML file updated successfully: %s Adding to git commit", path)
			//add file to git
			w, err := repo.Worktree()
			log.Printf("Adding file to git: %s", path)
			//strip /tmp/repo off filepath

			if _, err := os.Stat(path); os.IsNotExist(err) {
				log.Println("File does not exist:", path)
				// ... handle the error
			}
			relativePath := strings.TrimPrefix(path, "/tmp/repo/")
			_, err = w.Add(relativePath)
			if err != nil {
				log.Printf("Error adding file to git: %v", err)
				return
			}
			return
		}
	}
}

func (gs *Gitops) CheckForChangesAndCommit(repoConfig config.GitOps, repo *git.Repository, updateRequest map[string]string, runningData map[string]string) {
	// Check for changes and commit
	//log config available
	log.Println("Checking for changes and committing")
	w, err := repo.Worktree()
	status, err := w.Status()
	if err != nil {
		log.Printf("Error getting status: %v", err)
		return
	}
	log.Printf("Status: %v", status)
	if !status.IsClean() {
		// Commit the changes
		ntfyConfig := gs.cfg.Notifications.Ntfy
		ntfyManager := notifications.NewManager(ntfyConfig)
		message := fmt.Sprintf("🔔 *COMMIT!* 🔔\n\n*Container:* %s\n*Current Tag:* %s\n*New Tag:* %s\n*Found At:* %s", updateRequest["image"], runningData["image"], updateRequest["newTag"], updateRequest["repo"])
		payload := notifications.NotificationPayload{
			Message:  message,
			Priority: 1,
		}
		_, err = ntfyManager.SendNtfyNotificationGeneric(payload)
		// add files
		commit, err := w.Commit("Updated images", &git.CommitOptions{
			Author: &object.Signature{
				Name: "slackwatch",
				When: time.Now(),
			},
		})
		if err != nil {
			log.Printf("Error committing changes: %v", err)
			return
		}
		log.Printf("Commit: %v", commit)
		//push
		err = repo.Push(&git.PushOptions{
			Auth: &http.BasicAuth{
				Username: "slackwatch",
				Password: repoConfig.AuthToken,
			},
		})
		if err != nil {
			log.Printf("Error pushing changes: %v", err)
			return
		}
		log.Println("Changes pushed successfully to repository %s ", repoConfig.RepoURL)
	}
	//list files
	//commit
	//push
}

func (gs *Gitops) FindAndUpdateYaml(obj map[string]interface{}, currentImage string, currentTag string, newTag string) (map[string]interface{}, error) {
	// Get the image from the YAML
	//log config available
	log.Printf("Find and Update YAML: %v", obj)
	if containers, ok := obj["spec"].(map[interface{}]interface{})["template"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["containers"]; ok {
		for _, container := range containers.([]interface{}) {
			if container.(map[interface{}]interface{})["image"] != nil {
				image := container.(map[interface{}]interface{})["image"]
				if image == currentImage {
					log.Printf("FOUND IMAGE TO UPDATE: %s", image.(string))
					log.Print("Before: %v", container)
					newImage := strings.Replace(image.(string), currentTag, newTag, 1)
					log.Printf("New Image: %s", newImage)
					container.(map[interface{}]interface{})["image"] = newImage
					log.Printf("Updated Container: %v", container)
					log.Printf("Container image: %s", container.(map[interface{}]interface{})["image"])
					return obj, nil
				}

			}
		}

	}
	return nil, fmt.Errorf("image not found in YAML")
}

func (gs *Gitops) BuildPath(repo *git.Repository, updateRequest map[string]string) string {
	// Find the files to update
	//log config available
	log.Printf("Finding files to update for %s", updateRequest["repo"])
	if updateRequest["directory"] == "" {
		log.Printf("Subfolder specified: %s", updateRequest["directory"])
		subFolder := updateRequest["directory"]
		fullPath := "/tmp/repo/" + subFolder
		log.Printf("Full path: %s", fullPath)
		return fullPath
	} else {
		log.Printf("No subfolder specified, using name %s", updateRequest["name"])
		subFolder := updateRequest["name"]
		fullPath := "/tmp/repo/" + subFolder
		log.Printf("Full path: %s", fullPath)
		return fullPath

	}
	//list files

}

func (gs *Gitops) DeleteLocalRepo() {
	// Delete the cloned repository
	log.Println("Deleting the cloned repository")
	if err := os.RemoveAll("/tmp/repo"); err != nil {
		log.Printf("Error deleting repository: %v", err)
	}
	log.Println("Repository deleted successfully")

}

func (gs *Gitops) CloneRepo(repoConfig config.GitOps) (*git.Repository, error) {
	log.Printf("Cloning repository %s", repoConfig.RepoURL)
	// Clone the repository
	repo, err := git.PlainClone("/tmp/repo", false, &git.CloneOptions{
		URL:      repoConfig.RepoURL,
		Progress: log.Writer(),
		Auth: &http.BasicAuth{
			Username: "slackwatch",
			Password: repoConfig.AuthToken,
		},
	})
	if err != nil {
		log.Printf("Error cloning repository: %v", err)
		return nil, err
	}
	return repo, nil
}

// get repo config based on repo name
func (gs *Gitops) GetRepoConfig(repoName string) (config.GitOps, error) {
	// Get the configuration for the repository
	//log config available
	log.Printf("Getting repository configuration for %s", repoName)
	for _, repoConfig := range gs.cfg.GitOps {
		if repoConfig.Name == repoName {
			log.Printf("Repository config %s found for %s", repoConfig.Name, repoName)
			return repoConfig, nil
		}
	}
	// If no specific repository is matched, use the default repository
	return config.GitOps{}, fmt.Errorf("no repository configuration found for %s", repoName)
}

// compare to ensure tags are as expected
func (gs *Gitops) SanityChecks(updateRequest map[string]string, runningData map[string]string) bool {
	// Check if the new tag is the same as the current tag
	log.Printf("New Image vs Running Image: %s vs %s", updateRequest["image"], runningData["image"])
	if updateRequest["image"] != runningData["image"] {
		log.Printf("Update Request image %s does not match running image %s", updateRequest["image"], runningData["image"])
		return false
	}
	return true
}
