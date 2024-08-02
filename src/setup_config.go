package main

import (
	"context"
	"encoding/json"
	"fmt"
)

func setupConfig(ctx context.Context, rawConfig []byte) (*BigChange, error) {
	log := LoggerFromContext(ctx)
	bigChange := &BigChange{}

	err := json.Unmarshal(rawConfig, bigChange)
	if err != nil {
		log.Error("failed to unmarshal config", "error", err)
		return nil, err
	}

	if bigChange.Domains == nil {
		log.Error("missing or empty config field", "field", "BigChange.Domains")
		return nil, fmt.Errorf("missing or empty config field")
	}
	if bigChange.Settings == nil {
		log.Error("missing or empty config field", "field", "BigChange.Settings")
		return nil, fmt.Errorf("missing or empty config field")
	}

	if bigChange.Settings.MainBranch == "" {
		log.Error("missing or empty config field", "field", "BigChange.Settings.MainBranch")
		return nil, fmt.Errorf("missing or empty config field")
	}
	if bigChange.Settings.Remote == "" {
		log.Error("missing or empty config field", "field", "BigChange.Settings.Remote")
		return nil, fmt.Errorf("missing or empty config field")
	}
	if bigChange.Settings.BranchToSplit == "" {
		log.Error("missing or empty config field", "field", "BigChange.Settings.BranchToSplit")
		return nil, fmt.Errorf("missing or empty config field")
	}

	for _, domain := range bigChange.Domains {
		if domain.Path == "" {
			log.Error("missing or empty config field", "domain name", domain.Name, "field", "Domain.Path")
			return nil, fmt.Errorf("missing or empty config field")
		}
	}

	return bigChange, nil
}
