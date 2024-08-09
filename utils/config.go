package utils

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                int      `yaml:"lb_port"`
	Strategy            string   `yaml:"strategy"`
	Servers             []string `yaml:"servers"`
	HealthCheckInterval string   `yaml:"health_check_interval"`
	MaxRetryLimit       int      `yaml:"max_retry_limit"`
	ServerWeights       []int    `yaml:"server_weights"`
}

const MAX_LB_RETRY_LIMIT int = 3

func LoadConfig() (*Config, error) {
	var config Config

	yamlConfig, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return nil, err
	}

	if len(config.Servers) <= 0 {
		return nil, errors.New("config error: servers expected but not provided")
	}

	if config.Strategy == "round-robin" {
		weights := make([]int, len(config.Servers))
		for index := range weights {
			weights[index] = 1
		}
		config.ServerWeights = weights
	}

	if len(config.ServerWeights) <= 0 && config.Strategy == "weighted-round-robin" {
		return nil, errors.New("config error: servers weight are required for weighted algorithms")
	}

	if len(config.Servers) != len(config.ServerWeights) {
		return nil, errors.New("config error: number of servers does not match with number of server weight")
	}

	if config.Port == 0 {
		return nil, errors.New("config error: load balancer port expected but not provided")
	}

	// fmt.Print(config)

	return &config, nil
}
