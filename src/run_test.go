package main

import (
	"context"
	"fmt"
	"testing"
)

type givenRun struct {
	exportResults ExportResultsFunc
	flags         *Flags
	gitOps        *GitOps
	config        *BigChange
}

var runTests = []struct {
	description string
	given       givenRun
	expectedErr error
}{
	{
		description: "Happy path",
		given: givenRun{
			exportResults: checkExportResults(fixtureBigChange(func(bc *BigChange) {
				bc.Domains[0].Branch = &Branch{
					Name: "bit-dom1-big-change-split",
				}
				bc.Domains[0].PullRequest = &PullRequest{
					Title: "AA dom1: Big change split",
					Body:  "This change refers to this refactor for domain AA dom1: https://example.com",
					Url:   "bit-dom1-big-change-split/pr",
				}
				bc.Domains[1].Branch = &Branch{
					Name: "bit-dom2-big-change-split",
				}
				bc.Domains[1].PullRequest = &PullRequest{
					Title: "BB dom2: Big change split",
					Body:  "This change refers to this refactor for domain BB dom2: https://example.com",
					Url:   "bit-dom2-big-change-split/pr",
				}
				bc.Domains[2].Branch = &Branch{
					Name: "bit-dom3-big-change-split",
				}
				bc.Domains[2].PullRequest = &PullRequest{
					Title: "CC dom3: Big change split",
					Body:  "This change refers to this refactor for domain CC dom3: https://example.com",
					Url:   "",
				}
			}).Domains),
			flags:  fixtureFlags(),
			gitOps: fixtureGitOps(),
			config: fixtureBigChange(),
		},
	},
	{
		description: "Don't create branches and PRs on cleanup",
		given: givenRun{
			exportResults: checkExportResults(nil),
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
		description: "Fail on gitCheckout",
		given: givenRun{
			exportResults: checkExportResults(nil),
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
		given: givenRun{
			exportResults: checkExportResults(nil),
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
		given: givenRun{
			exportResults: checkExportResults(nil),
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
		given: givenRun{
			exportResults: checkExportResults(nil),
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
		given: givenRun{
			exportResults: checkExportResults(nil),
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
		given: givenRun{
			exportResults: checkExportResults(nil),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitCheckoutFiles = func(ctx context.Context, s1, s2 string) error { return fmt.Errorf("gitCheckoutFiles failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitCheckoutFiles failed"),
	},
	{
		description: "Fail on gitReset",
		given: givenRun{
			exportResults: checkExportResults(nil),
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
		given: givenRun{
			exportResults: checkExportResults(nil),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.gitPushSetUpstream = func(ctx context.Context, s1, s2 string) error { return fmt.Errorf("gitPushSetUpstream failed") }
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("gitPushSetUpstream failed"),
	},
	{
		description: "Fail on createPr",
		given: givenRun{
			exportResults: checkExportResults(nil),
			flags:         fixtureFlags(),
			gitOps: fixtureGitOps(func(g *GitOps) {
				g.createPr = func(ctx context.Context, s1 *Settings, s2, s3, s4 string) (string, error) {
					return "", fmt.Errorf("createPr failed")
				}
			}),
			config: fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("createPr failed"),
	},
	{
		description: "Fail on exportResults",
		given: givenRun{
			exportResults: func(ctx context.Context, f *Flags, bc *BigChange) error { return fmt.Errorf("exportResults failed") },
			flags:         fixtureFlags(),
			gitOps:        fixtureGitOps(),
			config:        fixtureBigChange(),
		},
		expectedErr: fmt.Errorf("exportResults failed"),
	},
}

func TestRun(t *testing.T) {
	ctxWithSilentLogger := ContextWithSilentLogger(context.Background())

	for _, tt := range runTests {
		t.Run(tt.description, func(t *testing.T) {
			bit := &BigIsTiny{
				exportResults: tt.given.exportResults,
				flags:         tt.given.flags,
				gitOps:        tt.given.gitOps,
			}
			gotErr := bit.run(ctxWithSilentLogger, tt.given.config)

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
