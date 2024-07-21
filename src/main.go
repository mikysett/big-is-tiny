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
	IsDraftPrs         bool
	RepoPath           string
	BranchNameTemplate string
	CommitMsgTemplate  string
	PrNameTemplate     string
	PrDescTemplate     string
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
	name        string
	description string
}

type ctxLogger struct{}

type Flags struct {
	Cleanup    bool
	Verbose    bool
	DryRun     bool
	ConfigPath string
	Platform   Platform
}

type GitZeroArgsFunc func(context.Context) error
type GitOneArgStringFunc func(context.Context, string) error
type GitTwoArgsStringFunc func(context.Context, string, string) error
type GitStatusFunc func(context.Context) ([]byte, error)

type BigIsTiny struct {
	chdirWithLogs func(context.Context, string) error
	flags         *Flags
	gitOps        *GitOps
}

type GitOps struct {
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
	flags, err := getFlags(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log := newLogger(flags.Verbose)
	ctx := ContextWithLogger(context.Background(), log)

	bigChange, err := setupConfig(ctx, flags.ConfigPath)
	if err != nil {
		os.Exit(1)
	}
	log.Debug("config extracted from config file", "bigChange", bigChange)

	bigIsTiny := BigIsTiny{
		chdirWithLogs: chdirWithLogs,
		flags:         flags,
		gitOps: &GitOps{
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
		},
	}

	err = bigIsTiny.run(ctx, bigChange)
	if err != nil {
		os.Exit(1)
	}
}

const usage = `Usage: bit [-v | --verbose] [-cleanup] [-dryrun] [-p | --platform] [-h | --help] <path to config file>

If not specified the default path to the config file is './bit_config.json'

  -cleanup
        delete branches and PRs
  -v, --verbose
        set logs to DEBUG level
  -dryrun
        do not create branches or PRs
  -p, --platform
		platform used for PRs, can be "github" (default) or "azure"
  -h, --help
        print this help information
`

func getFlags(progName string, args []string) (*Flags, error) {
	rawFlags := flag.NewFlagSet(progName, flag.ExitOnError)

	var verbose, cleanup, dryRun bool
	var rawPlatform string
	var platform Platform
	rawFlags.BoolVar(&cleanup, "cleanup", false, "delete branches and PRs")
	rawFlags.BoolVar(&verbose, "verbose", false, "set logs to DEBUG level")
	rawFlags.BoolVar(&verbose, "v", false, "set logs to DEBUG level")
	rawFlags.BoolVar(&dryRun, "dryrun", false, "do not create branches or PRs")
	rawFlags.StringVar(&rawPlatform, "platform", "github", "platform used for PRs, can be `github` (default) or `azure`")
	rawFlags.StringVar(&rawPlatform, "p", "github", "platform used for PRs, can be `github` (default) or `azure`")
	rawFlags.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	rawFlags.Parse(args)

	switch rawPlatform {
	case "github":
		platform = Platform(GitHub)
	case "azure":
		platform = Platform(Azure)
	default:
		return nil, fmt.Errorf("platform '%s' is not supported", rawPlatform)
	}

	flags := &Flags{
		Cleanup:  cleanup,
		Verbose:  verbose,
		DryRun:   dryRun,
		Platform: platform,
	}
	if configPath := rawFlags.Arg(0); configPath != "" {
		flags.ConfigPath = configPath
	} else {
		flags.ConfigPath = "bit_config.json"
	}

	return flags, nil
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
