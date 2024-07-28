package main

import (
	"context"
	"encoding/json"
)

type AzurePr struct {
	URL string `json:"url"`
}

func AzureCreatePr(ctx context.Context, settings *Settings, head, title, description string) (string, error) {
	prFlags := []string{
		"repos", "pr", "create",
		"--source-branch", head,
		"--title", title,
		"--description", description,
		"--target-branch", settings.MainBranch,
		"--output", "json",
	}
	if settings.IsDraftPrs {
		prFlags = append(prFlags, "--draft")
	}

	resp, err := runCmd(ctx, "az", prFlags...)
	if err != nil {
		return "", err
	}

	var pr AzurePr
	if err := json.Unmarshal(resp, &pr); err != nil {
		log := LoggerFromContext(ctx)
		log.Error("failed to unmarshal url of the PR", "error", err)
		return "", err
	}
	return string(pr.URL), nil
}
