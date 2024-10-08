package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func (bigChange *BigChange) generateFromTemplate(domain *Domain, template string) string {
	replacements := []string{
		"{{change_id}}", bigChange.Id,
		"{{domain_id}}", domain.Id,
		"{{domain_name}}", domain.Name,
		"{{pr_title}}", domain.PullRequest.Title,
		"{{pr_url}}", domain.PullRequest.Url,
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

func (bit *BigIsTiny) cleanup(ctx context.Context, bigChange *BigChange) {
	log := LoggerFromContext(ctx)
	log.Debug("remove all branches and PRs")

	_ = bit.gitOps.gitCheckout(ctx, bigChange.Settings.MainBranch)
	var wg sync.WaitGroup
	for _, domain := range bigChange.Domains {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Errors are expected to be logged here as branch existence is not checked
			_ = bit.gitOps.gitDeleteBranch(ctx, domain.Branch.Name)
			_ = bit.gitOps.gitDeleteRemoteBranch(ctx, bigChange.Settings.Remote, domain.Branch.Name)
			_ = bit.gitOps.abandonPr(ctx, domain.Branch.Name)
		}()
	}
	wg.Wait()
}
