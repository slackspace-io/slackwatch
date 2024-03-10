package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Kubernetes KubernetesConfig `yaml:"kubernetes"`
    Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
    Name     string `yaml:"name"`
    Token    string `yaml:"token,omitempty"`
    Username string `yaml:"username,omitempty"`
    Password string `yaml:"password,omitempty"`
}

type RepositoryConfig struct {
    URL      string `yaml:"url"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}

type KubernetesConfig struct {
    PollingInterval int `yaml:"pollingInterval"`
    UseInClusterConfig bool   `yaml:"useInClusterConfig"`
    OutOfClusterConfig struct {
        KubeconfigPath string `yaml:"kubeconfigPath"`
    } `yaml:"outOfClusterConfig"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(configPath string) (*Config, error) {
    //log conf file path
    log.Printf("Loading configuration from %s", configPath)
    //log contents of conf file
    content, readErr := ioutil.ReadFile(configPath)
    if readErr != nil {
        log.Fatalf("Error reading config file: %s", readErr)
    }
    log.Printf("Config file contents: \n%s\n", string(content))
    bytes, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    var config Config
    if err := yaml.Unmarshal(bytes, &config); err != nil {
        return nil, err
    }
    //log config
    log.Printf("Config: %+v", config)
    return &config, nil
}
