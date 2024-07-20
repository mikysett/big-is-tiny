package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type BigChange struct {
	Domains  []*Domain
	Settings *Settings
}

type Settings struct {
	MainBranch         string
	Remote             string
	BranchToSplit      string
	IsDryRun           bool
	IsDraftPrs         bool
	RepoPath           string
	BranchNameTemplate string
	CommitMsgTemplate  string
	PrNameTemplate     string
	PrDescription      string
	Platform           Platform
}

type Domain struct {
	Name        string
	Id          string
	Path        string
	Teams       []Team
	Branch      *Branch
	PullRequest *PullRequest
}

type Team struct {
	Name string
	Url  string
	Type Communication
}

type Communication int

const (
	Slack Communication = iota
	Teams
)

type Platform int

const (
	Azure Platform = iota
	GitHub
)

type Branch struct {
	name string
}

type PullRequest struct {
	name string
}

type ctxLogger struct{}

type Flags struct {
	Cleanup    bool
	Verbose    bool
	ConfigPath string
}

type GitZeroArgsFunc func(context.Context) error
type GitOneArgStringFunc func(context.Context, string) error
type GitTwoArgsStringFunc func(context.Context, string, string) error
type GitStatusFunc func(context.Context) ([]byte, error)

type BigIsTiny struct {
	flags                 Flags
	gitCheckout           GitOneArgStringFunc
	gitCheckoutNewBranch  GitOneArgStringFunc
	gitDeleteBranch       GitOneArgStringFunc
	gitDeleteRemoteBranch GitTwoArgsStringFunc
	gitStatus             GitStatusFunc
	gitAdd                GitOneArgStringFunc
	gitCommit             GitOneArgStringFunc
	gitCheckoutFiles      GitOneArgStringFunc
	gitReset              GitZeroArgsFunc
	gitPushSetUpstream    GitTwoArgsStringFunc
}

func main() {
	flags := getFlags()

	ctx := ContextWithLogger(context.Background(), newLogger(flags.Verbose))
	bigIsTiny := BigIsTiny{
		flags:                 flags,
		gitCheckout:           gitCheckout,
		gitCheckoutNewBranch:  gitCheckoutNewBranch,
		gitDeleteBranch:       gitDeleteBranch,
		gitDeleteRemoteBranch: gitDeleteRemoteBranch,
		gitStatus:             gitStatus,
		gitAdd:                gitAdd,
		gitCommit:             gitCommit,
		gitCheckoutFiles:      gitCheckoutFiles,
		gitReset:              gitReset,
		gitPushSetUpstream:    gitPushSetUpstream,
	}

	err := bigIsTiny.run(ctx)
	if err != nil {
		os.Exit(1)
	}
}

const usage = `Usage of big-is-tiny:
  -cleanup
        delete branches and PRs
  -v, --verbose
        set logs to DEBUG level
  -h, --help
        print this help information
`

func getFlags() Flags {
	var verbose bool
	var cleanup bool
	flag.BoolVar(&cleanup, "cleanup", false, "delete branches and PRs")
	flag.BoolVar(&verbose, "verbose", false, "set logs to DEBUG level")
	flag.BoolVar(&verbose, "v", false, "set logs to DEBUG level")
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	flags := Flags{
		Cleanup: cleanup,
		Verbose: verbose,
	}
	if configPath := flag.Arg(0); configPath != "" {
		flags.ConfigPath = configPath
	} else {
		flags.ConfigPath = "config.json"
	}

	return flags
}

func newLogger(verbose bool) *slog.Logger {
	// Default level is Info
	var programLevel = new(slog.LevelVar)

	if verbose {
		programLevel.Set(slog.LevelDebug)
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	structuredLog := slog.New(jsonHandler)

	return structuredLog
}

func ContextWithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, log)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*slog.Logger); ok {
		return l
	}
	return newLogger(false)
}
