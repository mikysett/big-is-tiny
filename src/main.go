package main

import (
	"fmt"
	"log/slog"
	"os"
)

type BigChange struct {
	Domains  []Domain
	Settings Settings
}

type Settings struct {
	IsDryRun   bool
	IsDraftPrs bool
	Verbose    bool
	RepoPath   string
	Platform   Platform
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

func main() {
	log := setLog()

	err := run("config.json", log)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func setLog() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	structuredLog := slog.New(jsonHandler)

	return structuredLog
}
