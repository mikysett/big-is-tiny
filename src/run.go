package main

import (
	"fmt"
	"log/slog"
)

func run(configPath string, log *slog.Logger) error {

	bigChange, err := setupConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to setup config from %s: %w", configPath, err)
	}
	log.Info("after configuration file", "bigChange", bigChange)
	return nil
}
