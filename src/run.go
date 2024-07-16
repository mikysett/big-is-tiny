package main

import (
	"fmt"
	"log/slog"
	"syscall"
)

func run(configPath string, log *slog.Logger) error {

	bigChange, err := setupConfig(configPath)
	if err != nil {
		log.Error("failed to setup config", "configPath", configPath, "error", err)
		return err
	}
	log.Info("config extracted from config file", "bigChange", bigChange)

	err = syscall.Chdir(bigChange.Settings.RepoPath)
	if err != nil {
		log.Error("failed to change directory", "target directory", bigChange.Settings.RepoPath, "error", err)
		return err
	}

	err = gitCheckout(bigChange.Settings.MainBranch)
	if err != nil {
		log.Error("failed to git checkout on main branch", "branch", bigChange.Settings.MainBranch, "error", err)
		return err
	}

	err = gitCheckoutFiles(bigChange.Settings.BranchToSplit)
	if err != nil {
		log.Error("failed to git checkout files from branch to split", "branch", bigChange.Settings.BranchToSplit, "error", err)
		return err
	}

	branches, err := createBranches(bigChange)
	if err != nil {
		log.Error("failed to create branches", "error", err)
		return err
	}

	pullRequests, err := createPullRequests(bigChange, branches)
	if err != nil {
		return err
	}

	log.Info("Pull Requests created", "pullRequests", pullRequests)
	return nil
}

func createBranches(bigChange *BigChange) ([]Branch, error) {
	return nil, fmt.Errorf("not implemented")
}

func createPullRequests(bigChange *BigChange, branches []Branch) ([]*PullRequest, error) {
	return nil, fmt.Errorf("not implemented")
}
