package catalog

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var ConfigDefault []byte

type Config struct {
	DBString string `yaml:"db_string"`
	Port     int    `yaml:"port"`
}

func NewConfig(path string) (*Config, error) {
	var (
		err         error
		configBytes = ConfigDefault
	)
	if path != "" {
		configBytes, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}
	var cfg Config
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}
