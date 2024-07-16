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
	log.Info("config extracted from config file", "bigChange", bigChange)

	// Check if the repo exists

	// Checkout to main

	// Checkout the files from the big change branch

	// Create a single branch and then a PR for each domain
	return nil
}
