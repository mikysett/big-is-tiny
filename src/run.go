package main

import (
	"context"
	"strings"
)

func (bit *BigIsTiny) run(ctx context.Context, config *BigChange) (err error) {
	// On failure remove the branches created during the split
	defer func() {
		if err != nil {
			bit.cleanup(ctx, config)
		}
	}()

	// All branches have need to be checked out from the main branch
	err = bit.gitOps.gitCheckout(ctx, config.Settings.MainBranch)
	if err != nil {
		return err
	}

	// We fetch the files from the change we are working on
	err = bit.gitOps.gitCheckoutFiles(ctx, config.Settings.BranchToSplit)
	if err != nil {
		return err
	}

	// Un-stage the files added from the big change branch
	err = bit.gitOps.gitReset(ctx)
	if err != nil {
		return err
	}

	// Get a list of all touched files
	changedFiles, err := bit.listChangedFiles(ctx)
	if err != nil {
		return err
	}
	log := LoggerFromContext(ctx)
	log.Debug("changed files", "files", changedFiles)

	for _, domain := range config.Domains {
		if !fileChangedInDomain(domain.Path, changedFiles) {
			continue
		}
		domain.Branch = &Branch{
			name: config.generateFromTemplate(domain, config.Settings.BranchNameTemplate),
		}
		domain.PullRequest = &PullRequest{
			title: config.generateFromTemplate(domain, config.Settings.PrNameTemplate),
			body:  config.generateFromTemplate(domain, config.Settings.PrDescTemplate),
		}

		if bit.flags.Cleanup {
			continue
		}

		err = bit.createBranch(ctx, config, domain, config.Settings)
		if err != nil {
			return err
		}

		domain.PullRequest.url, err = bit.createPullRequest(ctx, domain, config.Settings)
		if err != nil {
			return err
		}
	}

	if bit.flags.Cleanup {
		bit.cleanup(ctx, config)
	}

	return nil
}

func (bit *BigIsTiny) listChangedFiles(ctx context.Context) ([]string, error) {
	gitStatusResponse, err := bit.gitOps.gitStatus(ctx)
	if err != nil {
		return nil, err
	}

	rawOutput := string(gitStatusResponse[:])
	rawList := strings.Split(strings.Trim(rawOutput, "\n"), "\n")

	changedFiles := make([]string, 0, len(rawList))
	for _, statusLine := range rawList {
		_, filePath, _ := strings.Cut(strings.TrimSpace(statusLine), " ")
		changedFiles = append(changedFiles, strings.Trim(filePath, "\""))
	}
	return changedFiles, nil
}

func fileChangedInDomain(domainPath string, changedFiles []string) bool {
	for _, filePath := range changedFiles {
		if strings.HasPrefix(filePath, domainPath) {
			return true
		}
	}
	return false
}

func (bit *BigIsTiny) createBranch(ctx context.Context, config *BigChange, domain *Domain, settings *Settings) (err error) {
	defer func() {
		if err != nil {
			log := LoggerFromContext(ctx)
			log.Error("failed to create Branch", "branch", domain.Branch.name)
		}
	}()

	err = bit.gitOps.gitCheckoutNewBranch(ctx, domain.Branch.name)
	if err != nil {
		return err
	}

	err = bit.gitOps.gitAdd(ctx, domain.Path)
	if err != nil {
		return err
	}

	err = bit.gitOps.gitCommit(ctx, config.generateFromTemplate(domain, settings.CommitMsgTemplate))
	if err != nil {
		return err
	}

	err = bit.gitOps.gitPushSetUpstream(ctx, settings.Remote, domain.Branch.name)
	if err != nil {
		return err
	}

	// We go back to main branch not to change the repository initial state
	err = bit.gitOps.gitCheckout(ctx, settings.MainBranch)
	if err != nil {
		return err
	}

	return nil
}

func (bit *BigIsTiny) createPullRequest(ctx context.Context, domain *Domain, settings *Settings) (url string, err error) {
	url, err = bit.gitOps.createPr(ctx, settings, domain.Branch.name, domain.PullRequest.title, domain.PullRequest.body)
	if err != nil {
		log := LoggerFromContext(ctx)
		log.Error("failed to create Pull Request", "branch", domain.Branch.name)
		return "", err
	}

	return url, nil
}
