package config

import (
	"encoding/json"
	"fmt"
	"os"

	"devops/backend/internal/model"
)

func LoadHosts(path string) ([]model.Host, error) {
	if path == "" {
		return nil, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read hosts file: %w", err)
	}

	var hosts []model.Host
	if err := json.Unmarshal(raw, &hosts); err != nil {
		return nil, fmt.Errorf("parse hosts file: %w", err)
	}

	return hosts, nil
}
