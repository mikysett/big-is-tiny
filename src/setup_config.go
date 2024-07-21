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
			Remote:             "origin",
			BranchToSplit:      "big-change-to-split",
			IsDraftPrs:         false,
			RepoPath:           "../../test_repo",
			BranchNameTemplate: "bit-{{domain_name}}-big-change-split",
			CommitMsgTemplate:  "implement new feature for {{domain_name}} at {{team_name_1}}({{team_url_1}}) and {{team_name_2}}({{team_url_2}})",
			PrNameTemplate:     "{{domain_id}} {{domain_name}}: Big change split",
			PrDescTemplate:     "This change refers to this refactor for domain {{domain_id}} {{domain_name}}: https://example.com",
			Platform:           Platform(Azure),
		},
		Domains: []*Domain{
			{
				Name: "dom1",
				Id:   "AA",
				Path: "domains/dom1/",
				Teams: []Team{
					{
						Name: "First Team AA",
						Url:  "https://example_1.com",
						Type: Communication(Slack),
					},
				},
			},
			{
				Name: "dom2",
				Id:   "BB",
				Path: "domains/dom2/",
				Teams: []Team{
					{
						Name: "Team BB 1",
						Url:  "https://example_2.com",
						Type: Communication(Teams),
					},
					{
						Name: "Team BB 2",
						Url:  "https://example_2_bis.com",
						Type: Communication(Slack),
					},
				},
			},
		},
	}, nil
}
