package main

import (
	"context"
	"strings"
)

func GitHubCreatePr(ctx context.Context, settings *Settings, head, title, body string) (string, error) {
	prFlags := []string{
		"pr", "create",
		"-H", head,
		"-t", title,
		"-b", body,
		"-B", settings.MainBranch,
	}
	if settings.IsDraftPrs {
		prFlags = append(prFlags, "-d")
	}

	_, err := runCmd(ctx, "gh", prFlags...)
	if err != nil {
		return "", err
	}

	rawPrUrl, err := runCmd(ctx, "gh", "pr", "view", head, "--json", "url", "--template", "'{{.url}}'")
	if err != nil {
		return "", err
	}
	prUrl := string(rawPrUrl[:])
	return strings.Trim(prUrl, "'"), nil
}

// GitHub automatically abandon PR with deleted source branches, so this is a noOp
func GitHubAbandonPr(_ context.Context, _ string) error {
	return nil
}
