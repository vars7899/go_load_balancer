package utils

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type CoreConfig struct {
	StartComment string `yaml:"start_comment"`
}

func LoadCoreConfig() (*CoreConfig, error) {
	var config CoreConfig
	yaml_file, err := os.ReadFile("config.core.yaml")
	if err != nil {
		return nil, fmt.Errorf("core configuration error: %s", err.Error())
	}

	err = yaml.Unmarshal(yaml_file, &config)
	if err != nil {
		return nil, fmt.Errorf("core configuration error: %s", err.Error())
	}
	return &config, nil
}

func StartMsg(config *CoreConfig) {
	fmt.Printf("\n%v\n\n", config.StartComment)
}

var Lg *zap.Logger

func InitZapLogger() *zap.Logger {
	Lg, _ = zap.NewProduction()
	return Lg
}
