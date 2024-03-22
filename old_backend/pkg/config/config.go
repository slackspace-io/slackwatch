package config

import (
	"io/ioutil"
	"log"
	"os" // Added to import the os package

	"gopkg.in/yaml.v2"
)

type Config struct {
	Kubernetes    KubernetesConfig `yaml:"kubernetes"`
	Repositories  []Repository     `yaml:"repositories"`
	System        SystemConfig     `yaml:"system"`
	Magic         MagicConfig      `yaml:"magic"`
	GitOps        []GitOps         `yaml:"gitops"`
	Notifications struct {
		Ntfy NtfyConfig `yaml:"ntfy"`
	} `yaml:"notifications"`
}

type GitOps struct {
	Name      string `yaml:"name"`
	RepoURL   string `yaml:"repoURL"`
	Branch    string `yaml:"branch"`
	AuthToken string `yaml:"authToken"`
}

type Repository struct {
	Name        string `yaml:"name"`
	Username    string `yaml:"username,omitempty"`
	Password    string `yaml:"password,omitempty"`
	Token       string `yaml:"token,omitempty"`
	DefaultRepo bool   `yaml:"defaultRepo,omitempty"` // New field to indicate if this is the default repository
}

type KubernetesConfig struct {
	PollingInterval    int  `yaml:"pollingInterval"`
	UseInClusterConfig bool `yaml:"useInClusterConfig"`
	OutOfClusterConfig struct {
		KubeconfigPath string `yaml:"kubeconfigPath"`
	} `yaml:"outOfClusterConfig"`
}

type SystemConfig struct {
	Schedule     string `yaml:"schedule"`
	RunAtStartup bool   `yaml:"runAtStartup"`
}

type MagicConfig struct {
	ExcludePatterns []string `yaml:"excludePatterns"`
	IncludePatterns []string `yaml:"includePatterns"` // Add this line
}

type NtfyConfig struct {
	URL      string `yaml:"url"`
	Topic    string `yaml:"topic"`
	Token    string `yaml:"token"`
	Priority int    `yaml:"priority"`
	Reminder string `yaml:"reminder"` // Add this line
}

// LoadConfig reads and parses the configuration file
func LoadConfig(configPath string) (*Config, error) {
	log.Printf("Starting to load configuration from %s", configPath)

	// Read the configuration file
	content, readErr := ioutil.ReadFile(configPath)
	if readErr != nil {
		log.Printf("Error reading config file: %s", readErr)
		return nil, readErr
	}
	log.Printf("Successfully read config file: %s", configPath)

	// Unmarshal the YAML content into the Config struct
	var config Config
	if err := yaml.Unmarshal(content, &config); err != nil {
		log.Printf("Error unmarshalling config file: %s", err)
		return nil, err
	}

	// New code to read NTFY_TOKEN from environment variable
	ntfyToken := os.Getenv("NTFY_TOKEN")
	if ntfyToken != "" {
		config.Notifications.Ntfy.Token = ntfyToken
	}

	//get token for matching gitops repo config
	for i, gitops := range config.GitOps {
		if gitops.AuthToken == "" {
			gitopsToken := os.Getenv(gitops.Name + "-token")
			if gitopsToken != "" {
				log.Printf("Found token for %s", gitops.Name)
				config.GitOps[i].AuthToken = gitopsToken
			}
		}
	}

	log.Printf("Successfully loaded configuration: %+v", config)

	return &config, nil
}