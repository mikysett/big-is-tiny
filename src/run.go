package main

import (
	"context"
	"fmt"
	"strings"
)

func (bit *BigIsTiny) run(ctx context.Context) (err error) {
	log := LoggerFromContext(ctx)

	bigChange, err := setupConfig(ctx, bit.flags.ConfigPath)
	if err != nil {
		return err
	}
	log.Debug("config extracted from config file", "bigChange", bigChange)

	err = bit.chdirWithLogs(ctx, bigChange.Settings.RepoPath)
	if err != nil {
		return err
	}

	// TODO: add a defer here to cleanup branches, PRs and un-staged files in case of failure or success accordingly
	defer func() {
		if err != nil {
			bit.cleanup(ctx, bigChange)
		}
	}()

	// All branches have need to be checked out from the main branch
	err = bit.gitOps.gitCheckout(ctx, bigChange.Settings.MainBranch)
	if err != nil {
		return err
	}

	// We fetch the files from the change we are working on
	err = bit.gitOps.gitCheckoutFiles(ctx, bigChange.Settings.BranchToSplit)
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
	log.Debug("changed files", "files", changedFiles)

	for _, domain := range bigChange.Domains {
		if !fileChangedInDomain(domain.Path, changedFiles) {
			continue
		}
		domain.Branch = &Branch{
			name: generateFromTemplate(domain, bigChange.Settings.BranchNameTemplate),
		}

		if bit.flags.Cleanup {
			continue
		}

		err = bit.createBranch(ctx, domain, bigChange.Settings)
		if err != nil {
			return err
		}

		// TODO: create the pull request
		// domain.PullRequest, err = createPullRequest(ctx, domain, bigChange.Settings)
		// if err != nil {
		// 	return err
		// }
	}

	if bit.flags.Cleanup {
		bit.cleanup(ctx, bigChange)
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

func (bit *BigIsTiny) createBranch(ctx context.Context, domain *Domain, settings *Settings) error {
	err := bit.gitOps.gitCheckoutNewBranch(ctx, domain.Branch.name)
	if err != nil {
		return err
	}

	err = bit.gitOps.gitAdd(ctx, domain.Path)
	if err != nil {
		return err
	}

	err = bit.gitOps.gitCommit(ctx, generateFromTemplate(domain, settings.CommitMsgTemplate))
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

func generateFromTemplate(domain *Domain, template string) string {
	replacements := []string{
		"{{domain_name}}", domain.Name,
		"{{domain_id}}", domain.Id,
	}
	for i, team := range domain.Teams {
		replacements = append(replacements,
			fmt.Sprintf("{{team_name_%d}}", i+1), team.Name,
			fmt.Sprintf("{{team_url_%d}}", i+1), team.Url,
		)
	}
	r := strings.NewReplacer(replacements...)
	return r.Replace(template)
}

func createPullRequest(ctx context.Context, domain *Domain, settings *Settings) (pr *PullRequest, err error) {
	defer func() {
		if err != nil {
			log := LoggerFromContext(ctx)
			log.Error("failed to create Pull Request", "error", err)
		}
	}()
	return nil, fmt.Errorf("not implemented")
}

func (bit *BigIsTiny) cleanup(ctx context.Context, bigChange *BigChange) {
	log := LoggerFromContext(ctx)
	log.Info("remove all branches and PRs")

	_ = bit.gitOps.gitCheckout(ctx, bigChange.Settings.MainBranch)
	for _, domain := range bigChange.Domains {
		if domain.Branch == nil {
			continue
		}
		_ = bit.gitOps.gitDeleteBranch(ctx, domain.Branch.name)
		_ = bit.gitOps.gitDeleteRemoteBranch(ctx, bigChange.Settings.Remote, domain.Branch.name)
	}
}
