package config

import (
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Kubernetes   KubernetesConfig `yaml:"kubernetes"`
    Repositories []Repository     `yaml:"repositories"`
    System       SystemConfig     `yaml:"system"`
}

type Repository struct {
    Name        string `yaml:"name"`
    Username    string `yaml:"username,omitempty"`
    Password    string `yaml:"password,omitempty"`
    Token       string `yaml:"token,omitempty"`
    DefaultRepo bool   `yaml:"defaultRepo,omitempty"` // New field to indicate if this is the default repository
}

type KubernetesConfig struct {
    PollingInterval    int    `yaml:"pollingInterval"`
    UseInClusterConfig bool   `yaml:"useInClusterConfig"`
    OutOfClusterConfig struct {
        KubeconfigPath string `yaml:"kubeconfigPath"`
    } `yaml:"outOfClusterConfig"`
}

type SystemConfig struct {
    Schedule string `yaml:"schedule"`
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
    log.Printf("Successfully loaded configuration: %+v", config)

    return &config, nil
}
