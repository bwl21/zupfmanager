package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// configService implements ConfigService interface
type configService struct {
	defaultConfigPath string
}

// NewConfigService creates a new config service
func NewConfigService() ConfigService {
	return &configService{
		defaultConfigPath: "default-project-config.json",
	}
}

// NewConfigServiceWithPath creates a new config service with custom default path
func NewConfigServiceWithPath(defaultPath string) ConfigService {
	return &configService{
		defaultConfigPath: defaultPath,
	}
}

// LoadFromFile loads configuration from a file
func (c *configService) LoadFromFile(path string) (map[string]interface{}, error) {
	if path == "" {
		return map[string]interface{}{}, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON from %s: %w", path, err)
	}

	return config, nil
}

// LoadDefault loads the default configuration
func (c *configService) LoadDefault() (map[string]interface{}, error) {
	return c.LoadFromFile(c.defaultConfigPath)
}
