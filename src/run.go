package main

import (
	"context"
	"strings"

	"golang.org/x/sync/errgroup"
)

func (bit *BigIsTiny) run(ctx context.Context, config *BigChange) (err error) {
	// Generate the names for all new branches and PRs
	for _, domain := range config.Domains {
		domain.initDomain(config)
		// domain.Branch = &Branch{
		// 	Name: config.generateFromTemplate(domain, config.Settings.BranchNameTemplate),
		// }
		// domain.PullRequest = &PullRequest{
		// 	Title: config.generateFromTemplate(domain, config.Settings.PrNameTemplate),
		// 	Body:  config.generateFromTemplate(domain, config.Settings.PrDescTemplate),
		// }
	}

	// On cleanup or failure remove the branches and PRs created during the split
	defer func() {
		if bit.flags.Cleanup || err != nil {
			bit.cleanup(ctx, config)
		}
	}()

	if bit.flags.Cleanup {
		return nil
	}

	// Checkout to the main branch
	err = bit.gitOps.gitCheckout(ctx, config.Settings.MainBranch)
	if err != nil {
		return err
	}

	// Fetch the files on remote branch for the changes we are working on
	err = bit.gitOps.gitCheckoutFiles(ctx, config.Settings.Remote, config.Settings.BranchToSplit, bit.flags.AllowDeletions)
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

	errGrp := new(errgroup.Group)
	for _, domain := range config.Domains {
		if !fileChangedInDomain(domain.Path, changedFiles) {
			continue
		}

		err = bit.createBranch(ctx, config, domain, config.Settings)
		if err != nil {
			return err
		}

		errGrp.Go(func() error {
			err = bit.gitOps.gitPushSetUpstream(ctx, config.Settings.Remote, domain.Branch.Name)
			if err != nil {
				return err
			}

			domain.PullRequest.Url, err = bit.createPullRequest(ctx, domain, config.Settings)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := errGrp.Wait(); err != nil {
		return err
	}

	err = bit.exportResults(ctx, bit.flags, config)
	if err != nil {
		return err
	}

	return nil
}

func (domain *Domain) initDomain(config *BigChange) {
	domain.Branch = &Branch{
		Name: config.generateFromTemplate(domain, config.Settings.BranchNameTemplate),
	}
	domain.PullRequest = PullRequest{
		Title: config.generateFromTemplate(domain, config.Settings.PrNameTemplate),
		Body:  config.generateFromTemplate(domain, config.Settings.PrDescTemplate),
	}
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
			log.Error("failed to create Branch", "branch", domain.Branch.Name)
		}
	}()

	err = bit.gitOps.gitCheckoutNewBranch(ctx, domain.Branch.Name)
	if err != nil {
		return err
	}

	// We go back to main branch not to change the repository initial state
	defer func() {
		checkoutErr := bit.gitOps.gitCheckout(ctx, settings.MainBranch)
		if checkoutErr != nil {
			err = checkoutErr
		}
	}()

	err = bit.gitOps.gitAdd(ctx, domain.Path)
	if err != nil {
		return err
	}

	err = bit.gitOps.gitCommit(ctx, config.generateFromTemplate(domain, settings.CommitMsgTemplate))
	if err != nil {
		return err
	}

	return nil
}

func (bit *BigIsTiny) createPullRequest(ctx context.Context, domain *Domain, settings *Settings) (url string, err error) {
	url, err = bit.gitOps.createPr(ctx, settings, domain.Branch.Name, domain.PullRequest.Title, domain.PullRequest.Body)
	if err != nil {
		log := LoggerFromContext(ctx)
		log.Error("failed to create Pull Request", "branch", domain.Branch.Name)
		return "", err
	}

	return url, nil
}
