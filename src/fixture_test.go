package main

import "context"

func fixtureFlags(mods ...func(*Flags)) *Flags {
	flags := &Flags{
		Cleanup:    false,
		Verbose:    false,
		DryRun:     false,
		ConfigPath: "bit_config.json",
		Platform:   Platform(GitHub),
	}
	for _, mod := range mods {
		mod(flags)
	}
	return flags
}

func fixtureBigChange(mods ...func(*BigChange)) *BigChange {
	flags := &BigChange{
		Settings: &Settings{
			MainBranch:         "main",
			Remote:             "origin",
			BranchToSplit:      "big-change-to-split",
			IsDraftPrs:         false,
			BranchNameTemplate: "bit-{{domain_name}}-big-change-split",
			CommitMsgTemplate:  "implement new feature for {{domain_name}} at {{team_name_1}}({{team_url_1}}) and {{team_name_2}}({{team_url_2}})",
			PrNameTemplate:     "{{domain_id}} {{domain_name}}: Big change split",
			PrDescTemplate:     "This change refers to this refactor for domain {{domain_id}} {{domain_name}}: https://example.com",
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
					},
					{
						Name: "Team BB 2",
						Url:  "https://example_2_bis.com",
					},
				},
			},
			{
				Name: "dom3",
				Id:   "CC",
				Path: "domains/dom3/",
				Teams: []Team{
					{
						Name: "Team CC 1",
						Url:  "https://example_2.com",
					},
					{
						Name: "Team CC 2",
						Url:  "https://example_2_bis.com",
					},
				},
			},
		},
	}
	for _, mod := range mods {
		mod(flags)
	}
	return flags
}

func fixtureGitOps(mods ...func(*GitOps)) *GitOps {
	gitOps := &GitOps{
		gitCheckout:           func(ctx context.Context, s string) error { return nil },
		gitCheckoutNewBranch:  func(ctx context.Context, s string) error { return nil },
		gitDeleteBranch:       func(ctx context.Context, s string) error { return nil },
		gitDeleteRemoteBranch: func(ctx context.Context, s1, s2 string) error { return nil },
		// Used to check if branches need to be created because files were modified
		gitStatus: func(ctx context.Context) ([]byte, error) {
			return []byte(" M domains/dom1/file1\n A domains/dom2/file2\n"), nil
		},
		gitAdd:             func(ctx context.Context, s string) error { return nil },
		gitCommit:          func(ctx context.Context, s string) error { return nil },
		gitCheckoutFiles:   func(ctx context.Context, s1, s2 string) error { return nil },
		gitReset:           func(ctx context.Context) error { return nil },
		gitPushSetUpstream: func(ctx context.Context, s1, s2 string) error { return nil },
		createPr: func(ctx context.Context, s1 *Settings, s2, s3, s4 string) (string, error) {
			return s2 + "/pr", nil
		},
	}
	for _, mod := range mods {
		mod(gitOps)
	}
	return gitOps
}
