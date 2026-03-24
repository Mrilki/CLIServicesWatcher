package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

func GetDefaultConf() *models.Config {
	return models.GetDefaultConf()
}

func Load(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Config file %s not found\nUse default config", path)
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: path=%s: %v", ErrRead, path, err)
	}

	var cfg models.Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: path=%s: %v", ErrParse, path, err)
	}
	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: path=%s: %v", ErrValidate, path, err)
	}
	fmt.Printf("Loaded config file %s\nTargets found: %d\n", path, len(cfg.Targets))
	return &cfg, nil
}
