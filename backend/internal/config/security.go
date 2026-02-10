package config

import (
	"fmt"
	"strings"
)

type SecurityConfig struct {
	APIKeys map[string]string
}

func ParseSecurity(apiKeysRaw string) (SecurityConfig, error) {
	cfg := SecurityConfig{APIKeys: map[string]string{}}
	if strings.TrimSpace(apiKeysRaw) == "" {
		return cfg, nil
	}

	entries := strings.Split(apiKeysRaw, ",")
	for _, entry := range entries {
		parts := strings.Split(strings.TrimSpace(entry), ":")
		if len(parts) != 2 {
			return cfg, fmt.Errorf("invalid API_KEYS entry %q", entry)
		}
		token := strings.TrimSpace(parts[0])
		role := strings.TrimSpace(parts[1])
		if token == "" || role == "" {
			return cfg, fmt.Errorf("invalid API_KEYS entry %q", entry)
		}
		cfg.APIKeys[token] = role
	}

	return cfg, nil
}
