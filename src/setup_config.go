package main

import "context"

func setupConfig(ctx context.Context, configPath string) (*BigChange, error) {
	// TODO: Replace hardcoded values with fetched ones

	// Log errors accordingly when values will be properly fetched
	// log := LoggerFromContext(ctx)
	// log.Error("failed to setup config", "configPath", configPath, "error", err)
	return &BigChange{
		Settings: &Settings{
			MainBranch:         "main",
			BranchToSplit:      "big-change-to-split",
			IsDryRun:           true,
			IsDraftPrs:         false,
			Verbose:            true,
			RepoPath:           "../../test_repo",
			BranchNameTemplate: "bit-{{domain_name}}-big-change-split",
			CommitMsgTemplate:  "implement new feature for {{domain_name}} at {{team_name_1}}({{team_url_1}}) and {{team_name_2}}({{team_url_2}})",
			PrNameTemplate:     "{{domain_id}} {{domain_name}}: Big change split",
			PrDescription:      "This change refers to this refactor: https://example.com",
			Platform:           Platform(Azure),
		},
		Domains: []*Domain{
			{
				Name: "dom1",
				Id:   "AA",
				Path: "domains/dom1",
				Teams: []Team{
					{
						Name: "First Team AA",
						Url:  "team_url1",
						Type: Communication(Slack),
					},
				},
			},
			{
				Name: "dom2",
				Id:   "BB",
				Path: "domains/dom2",
				Teams: []Team{
					{
						Name: "Team BB 1",
						Url:  "team_url2",
						Type: Communication(Teams),
					},
					{
						Name: "Team BB 2",
						Url:  "team_url2bis",
						Type: Communication(Slack),
					},
				},
			},
		},
	}, nil
}
