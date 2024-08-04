package main

import (
	"context"
	"encoding/json"
	"strconv"
)

type AzurePr struct {
	BaseUrl      string `json:"baseUrl"`
	CodeReviewId int    `json:"codeReviewId"`
}

func AzureCreatePr(ctx context.Context, settings *Settings, head, title, description string) (string, error) {
	prFlags := []string{
		"repos", "pr", "create",
		"--source-branch", head,
		"--title", title,
		"--description", description,
		"--target-branch", settings.MainBranch,
		"--output", "json",
		"--query", "{baseUrl:repository.webUrl, codeReviewId:codeReviewId}",
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

	fullUrl := string(pr.BaseUrl) + "/pullrequest/" + strconv.Itoa(pr.CodeReviewId)
	return fullUrl, nil
}

func AzureAbandonPr(ctx context.Context, sourceBranch string) error {
	resp, err := runCmd(ctx, "az", "repos", "pr", "list",
		"--top", "1",
		"--status", "active",
		"--source-branch", sourceBranch,
		"--output", "json",
		"--query", "[].{codeReviewId:codeReviewId}")
	if err != nil {
		return err
	}

	var activePrsOnSourceBranch []AzurePr
	if err := json.Unmarshal(resp, &activePrsOnSourceBranch); err != nil {
		log := LoggerFromContext(ctx)
		log.Error("failed to unmarshal the PR id", "error", err)
		return err
	}

	if len(activePrsOnSourceBranch) < 1 {
		return nil
	}

	_, err = runCmd(ctx, "az", "repos", "pr", "update", "--id", strconv.Itoa(activePrsOnSourceBranch[0].CodeReviewId), "--status", "abandoned")
	if err != nil {
		return err
	}

	return nil
}
