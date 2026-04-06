package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

func Load(path string, log *slog.Logger) (*models.Config, error) {
	log.Debug("Reading config file", "path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: path=%s: %v", ErrRead, path, err)
	}

	log.Debug("Parsing config file", "path", path)
	var cfg models.Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: path=%s: %v", ErrParse, path, err)
	}

	log.Debug("Validating config file", "path", path)
	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: path=%s: %v", ErrValidate, path, err)
	}
	log.Debug("Normalizing config")
	cfg.Normalize()

	log.Info("Config loaded successfully",
		"path", path,
		"targets", len(cfg.Targets),
		"default_timeout", cfg.Timeout,
	)
	return &cfg, nil
}
