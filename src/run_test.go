package main

import (
	"context"
	"fmt"
	"testing"
)

type given struct {
	chdirWithLogs func(context.Context, string) error
	flags         *Flags
	gitOps        *GitOps
	config        *BigChange
}

var runTests = []struct {
	description string
	given       given
	expectedErr error
}{
	{
		description: "Happy path",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps:        fixtureGitOps(),
			config:        fixtureBigChange(),
		},
	},
	{
		description: "Don't create branches and PRs on cleanup",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags: fixtureFlags(func(f *Flags) {
				f.Cleanup = true
			}),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitAdd = func(ctx context.Context, s string) error { return fmt.Errorf("gitAdd should not be called") }
				g.gitCommit = func(ctx context.Context, s string) error { return fmt.Errorf("gitCommit should not be called") }
				g.gitPushSetUpstream = func(ctx context.Context, s1, s2 string) error {
					return fmt.Errorf("gitPushSetUpstream should not be called")
				}
			}),
			config: fixtureBigChange(),
		},
	},
	{
		description: "Fail on change directory",
		given: given{
			chdirWithLogs: func(ctx context.Context, s string) error { return fmt.Errorf("chdir failed") },
			flags:         fixtureFlags(),
			gitOps:        fixtureGitOps(),
			config:        fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("chdir failed"),
	},
	{
		description: "Fail on gitCheckout",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitCheckout = func(ctx context.Context, s string) error { return fmt.Errorf("gitCheckout failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitCheckout failed"),
	},
	{
		description: "Fail on gitCheckoutNewBranch",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitCheckoutNewBranch = func(ctx context.Context, s string) error { return fmt.Errorf("gitCheckoutNewBranch failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitCheckoutNewBranch failed"),
	},
	{
		description: "Fail on gitStatus",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitStatus = func(ctx context.Context) ([]byte, error) { return nil, fmt.Errorf("gitStatus failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitStatus failed"),
	},
	{
		description: "Fail on gitAdd",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitAdd = func(ctx context.Context, s string) error { return fmt.Errorf("gitAdd failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitAdd failed"),
	},
	{
		description: "Fail on gitCommit",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitCommit = func(ctx context.Context, s string) error { return fmt.Errorf("gitCommit failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitCommit failed"),
	},
	{
		description: "Fail on gitCheckoutFiles",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitCheckoutFiles = func(ctx context.Context, s string) error { return fmt.Errorf("gitCheckoutFiles failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitCheckoutFiles failed"),
	},
	{
		description: "Fail on gitReset",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitReset = func(ctx context.Context) error { return fmt.Errorf("gitReset failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitReset failed"),
	},
	{
		description: "Fail on gitPushSetUpstream",
		given: given{
			chdirWithLogs: fixtureChdirWithLogs(),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitPushSetUpstream = func(ctx context.Context, s1, s2 string) error { return fmt.Errorf("gitPushSetUpstream failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitPushSetUpstream failed"),
	},
}

func TestRun(t *testing.T) {
	for _, tt := range runTests {
		t.Run(tt.description, func(t *testing.T) {
			bit := &BigIsTiny{
				chdirWithLogs: tt.given.chdirWithLogs,
				flags:         *tt.given.flags,
				gitOps:        tt.given.gitOps,
			}
			gotErr := bit.run(context.Background(), tt.given.config)

			// We get an error when we don't expect it or we don't get one when we expect it
			if tt.expectedErr != nil != (gotErr != nil) {
				t.Errorf("got '%v', want '%v'", gotErr, tt.expectedErr)
			}
			// We get a different error of what's expected
			if tt.expectedErr != nil && gotErr != nil &&
				tt.expectedErr.Error() != gotErr.Error() {
				t.Errorf("got '%v', want '%v'", gotErr, tt.expectedErr)
			}
		})
	}
}
