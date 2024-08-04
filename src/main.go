package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type BigChange struct {
	Id       string    `json:"id"`
	Domains  []*Domain `json:"domains"`
	Settings *Settings `json:"settings"`
}

type Settings struct {
	MainBranch         string `json:"mainBranch"`
	Remote             string `json:"remote"`
	BranchToSplit      string `json:"branchToSplit"`
	IsDraftPrs         bool   `json:"isDraftPrs"`
	BranchNameTemplate string `json:"branchNameTemplate"`
	CommitMsgTemplate  string `json:"commitMsgTemplate"`
	PrNameTemplate     string `json:"prNameTemplate"`
	PrDescTemplate     string `json:"prDescTemplate"`
}

type Domain struct {
	Name        string       `json:"name"`
	Id          string       `json:"id"`
	Path        string       `json:"path"`
	Teams       []Team       `json:"teams"`
	Branch      *Branch      `json:"branch"`
	PullRequest *PullRequest `json:"pullRequest"`
}

type Team struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Branch struct {
	Name string `json:"name"`
}

type PullRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Url   string `json:"url"`
}

type ctxLogger struct{}

type ExportResultsFunc func(context.Context, *Flags, *BigChange) error

type BigIsTiny struct {
	flags         *Flags
	exportResults ExportResultsFunc
	gitOps        *GitOps
}

type Flags struct {
	Cleanup     bool
	Verbose     bool
	ConfigPath  string
	Platform    Platform
	MarkdownOut bool
	FileOut     string
}

type GitZeroArgsFunc func(context.Context) error
type GitOneArgStringFunc func(context.Context, string) error
type GitTwoArgsStringFunc func(context.Context, string, string) error
type GitStatusFunc func(context.Context) ([]byte, error)
type CreatePrFunc func(context.Context, *Settings, string, string, string) (string, error)
type AbandonPrFunc func(context.Context, string) error

type GitOps struct {
	gitCheckout           GitOneArgStringFunc
	gitCheckoutNewBranch  GitOneArgStringFunc
	gitDeleteBranch       GitOneArgStringFunc
	gitDeleteRemoteBranch GitTwoArgsStringFunc
	gitStatus             GitStatusFunc
	gitAdd                GitOneArgStringFunc
	gitCommit             GitOneArgStringFunc
	gitCheckoutFiles      GitTwoArgsStringFunc
	gitReset              GitZeroArgsFunc
	gitPushSetUpstream    GitTwoArgsStringFunc
	createPr              CreatePrFunc
	abandonPr             AbandonPrFunc
}

func main() {
	flags, err := getFlags(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log := newLogger(flags.Verbose)
	ctx := ContextWithLogger(context.Background(), log)

	jsonConfig, err := os.ReadFile(flags.ConfigPath)
	if err != nil {
		log.Error("failed to read config file",
			"config file path", flags.ConfigPath,
			"error", err)
		os.Exit(2)
	}

	bigChange, err := setupConfig(ctx, jsonConfig)
	if err != nil {
		os.Exit(3)
	}
	log.Debug("config extracted from config file", "bigChange", bigChange)

	bigIsTiny := BigIsTiny{
		flags:         flags,
		exportResults: exportResults,
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
			createPr:              GetCreatePrForPlatform(flags.Platform),
			abandonPr:             GetAbandonPrForPlatform(flags.Platform),
		},
	}

	err = bigIsTiny.run(ctx, bigChange)
	if err != nil {
		os.Exit(4)
	}
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

// To reduce noise in unit tests
func ContextWithSilentLogger(ctx context.Context) context.Context {
	dummyHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.Level(slog.LevelError + 1)})
	dummyLogger := slog.New(dummyHandler)

	return context.WithValue(ctx, ctxLogger{}, dummyLogger)
}
