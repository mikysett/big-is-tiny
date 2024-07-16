package main

func setupConfig(configPath string) (*BigChange, error) {
	// TODO: Replace hardcoded values with fetched ones
	return &BigChange{
		Settings: Settings{
			MainBranch:    "main",
			BranchToSplit: "big-change-to-split",
			IsDryRun:      true,
			IsDraftPrs:    false,
			Verbose:       true,
			RepoPath:      "../../test_repo",
			Platform:      Platform(Azure),
		},
		Domains: []Domain{
			{
				Name: "dom1",
				Path: "domains/dom1",
				Teams: []Team{
					{
						TeamUrl:  "team_url1",
						TeamType: Communication(Slack),
					},
				},
			},
			{
				Name: "dom2",
				Path: "domains/dom2",
				Teams: []Team{
					{
						TeamUrl:  "team_url2",
						TeamType: Communication(Teams),
					},
					{
						TeamUrl:  "team_url2bis",
						TeamType: Communication(Slack),
					},
				},
			},
		},
	}, nil
}
