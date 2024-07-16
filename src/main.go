package main

import (
	"log/slog"
	"os"
)

type BigChange struct {
	Domains  []Domain
	Settings Settings
}

type Settings struct {
	MainBranch    string
	BranchToSplit string
	IsDryRun      bool
	IsDraftPrs    bool
	Verbose       bool
	RepoPath      string
	Platform      Platform
}

type Domain struct {
	Name  string
	Path  string
	Teams []Team
}

type Team struct {
	TeamUrl  string
	TeamType Communication
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

func main() {
	log := setLog()

	err := run("config.json", log)
	if err != nil {
		os.Exit(1)
	}
}

func setLog() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	structuredLog := slog.New(jsonHandler)

	return structuredLog
}
