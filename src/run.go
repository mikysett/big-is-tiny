package main

import (
	"context"
	"fmt"
	"strings"
	"syscall"
)

func run(ctx context.Context, configPath string) (err error) {
	log := LoggerFromContext(ctx)

	bigChange, err := setupConfig(ctx, configPath)
	if err != nil {
		return err
	}
	log.Info("config extracted from config file", "bigChange", bigChange)

	err = syscall.Chdir(bigChange.Settings.RepoPath)
	if err != nil {
		log.Error("failed to change directory", "target directory", bigChange.Settings.RepoPath, "error", err)
		return err
	}

	// TODO: add a defer here to cleanup branches, PRs and un-staged files in case of failure or success accordingly

	// All branches have need to be checked out from the main branch
	err = gitCheckout(ctx, bigChange.Settings.MainBranch)
	if err != nil {
		return err
	}

	// We fetch the files from the change we are working on
	err = gitCheckoutFiles(ctx, bigChange.Settings.BranchToSplit)
	if err != nil {
		return err
	}

	// Un-stage the files added from the big change branch
	err = gitReset(ctx)
	if err != nil {
		return err
	}

	// Get a list of all touched files
	changedFiles, err := listChangedFiles(ctx)
	if err != nil {
		return err
	}
	log.Info("changed files", "files", changedFiles)

	for _, domain := range bigChange.Domains {
		if !fileChangedInDomain(domain.Path, changedFiles) {
			continue
		}
		domain.Branch, err = createBranch(ctx, domain, bigChange.Settings)
		if err != nil {
			return err
		}

		// TODO: create the pull request
	}

	return nil
}

func listChangedFiles(ctx context.Context) ([]string, error) {
	gitStatusResponse, err := gitStatus(ctx)
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

func createBranch(ctx context.Context, domain *Domain, settings *Settings) (branch *Branch, err error) {
	name := generateFromTemplate(domain, settings.BranchNameTemplate)

	err = gitCheckoutNewBranch(ctx, name)
	if err != nil {
		return nil, err
	}

	// If something goes wrong we want to delete the branch to reduce garbage
	defer func() {
		if err != nil {
			_ = gitCheckout(ctx, settings.MainBranch)
			_ = gitDeleteBranch(ctx, name)
		}
	}()

	err = gitAdd(ctx, domain.Path)
	if err != nil {
		return nil, err
	}

	err = gitCommit(ctx, generateFromTemplate(domain, settings.CommitMsgTemplate))
	if err != nil {
		return nil, err
	}

	// TODO: git push set upstream

	// We go back to main branch not to change the repository initial state
	err = gitCheckout(ctx, settings.MainBranch)
	if err != nil {
		return nil, err
	}

	return &Branch{
		name: name,
	}, nil

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

func createPullRequest(ctx context.Context, bigChange *BigChange, branches []Branch) (*PullRequest, error) {
	return nil, fmt.Errorf("not implemented")
}
