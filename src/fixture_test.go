package main

import "context"

func fixtureChdirWithLogs() func(ctx context.Context, path string) error {
	return func(ctx context.Context, s string) error { return nil }
}

func fixtureFlags(mods ...func(*Flags)) *Flags {
	flags := &Flags{
		Cleanup:    false,
		Verbose:    false,
		ConfigPath: "",
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
		gitCheckoutFiles:   func(ctx context.Context, s string) error { return nil },
		gitReset:           func(ctx context.Context) error { return nil },
		gitPushSetUpstream: func(ctx context.Context, s1, s2 string) error { return nil },
	}
	for _, mod := range mods {
		mod(gitOps)
	}
	return gitOps
}
