package main

import (
	"context"
	"fmt"
	"strings"
)

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
