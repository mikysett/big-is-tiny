package main

import "context"

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

	resp, err := runCmd(ctx, "gh", prFlags...)

	if err != nil {
		return "", err
	}
	return string(resp[:]), nil
}
