package main

import (
	"context"
	"flag"
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
	Verbose            bool
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

func main() {
	flags := getFlags()

	ctx := ContextWithLogger(context.Background(), newLogger(flags.Verbose))
	err := run(ctx, flags)
	if err != nil {
		os.Exit(1)
	}
}

func getFlags() Flags {
	cleanup := flag.Bool("cleanup", false, "delete branches and PRs")
	verbose := flag.Bool("verbose", false, "set logs to DEBUG level")
	flag.Parse()

	flags := Flags{
		Cleanup: *cleanup,
		Verbose: *verbose,
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
