package main

import (
	"context"
	"fmt"
)

func gitCheckout(ctx context.Context, branchName string) error {
	_, err := runCmd(ctx, "git", "checkout", branchName)
	if err != nil {
		return err
	}
	return nil
}

func gitCheckoutNewBranch(ctx context.Context, branchName string) error {
	_, err := runCmd(ctx, "git", "checkout", "-b", branchName)
	if err != nil {
		return err
	}
	return nil
}

func gitDeleteBranch(ctx context.Context, branchName string) error {
	_, err := runCmd(ctx, "git", "branch", "-D", branchName)
	if err != nil {
		return err
	}
	return nil
}

func gitDeleteRemoteBranch(ctx context.Context, remote string, branchName string) error {
	_, err := runCmd(ctx, "git", "push", remote, "-d", branchName)
	if err != nil {
		return err
	}
	return nil
}

func gitStatus(ctx context.Context) ([]byte, error) {
	resp, err := runCmd(ctx, "git", "status", "--porcelain")
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func gitAdd(ctx context.Context, pathToAdd string) error {
	_, err := runCmd(ctx, "git", "add", pathToAdd)
	if err != nil {
		return err
	}
	return nil
}

func gitCommit(ctx context.Context, message string) error {
	_, err := runCmd(ctx, "git", "commit", "-m", message)
	if err != nil {
		return err
	}
	return nil
}

func gitCheckoutFiles(ctx context.Context, remote string, branchName string) error {
	_, err := runCmd(ctx, "git", "checkout", "--no-overlay", fmt.Sprintf("%s/%s", remote, branchName), "--", ".")
	if err != nil {
		return err
	}
	return nil
}

func gitReset(ctx context.Context) error {
	_, err := runCmd(ctx, "git", "reset")
	if err != nil {
		return err
	}
	return nil
}

func gitPushSetUpstream(ctx context.Context, remote string, branchName string) error {
	_, err := runCmd(ctx, "git", "push", "--set-upstream", remote, branchName)
	if err != nil {
		return err
	}
	return nil
}
