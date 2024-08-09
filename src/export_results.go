package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type createdPr struct {
	Branch string `json:"branch"`
	PrUrl  string `json:"prUrl"`
}

func exportResults(ctx context.Context, flags *Flags, config *BigChange) (err error) {
	var fdOut *os.File
	if flags.FileOut == "" {
		fdOut = os.Stdout
	} else {
		fdOut, err = os.OpenFile(flags.FileOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log := LoggerFromContext(ctx)
			log.Error("failed to create results file", "path", flags.FileOut, "error", err)
			return err
		}
		defer fdOut.Close()
	}

	createdPrs := make([]createdPr, 0, len(config.Domains))
	for _, domain := range config.Domains {
		if domain.PullRequest == nil || domain.PullRequest.Url == "" {
			continue
		}
		createdPrs = append(createdPrs, createdPr{
			Branch: domain.Branch.Name,
			PrUrl:  domain.PullRequest.Url,
		})
	}

	if flags.MarkdownOut {
		for _, pr := range createdPrs {
			fmt.Fprintf(fdOut, "[%s](%s)\n", pr.Branch, pr.PrUrl)
		}
	} else {
		jsonFormattedResult, err := json.MarshalIndent(createdPrs, "", "    ")
		if err != nil {
			log := LoggerFromContext(ctx)
			log.Error("failed to marshal results", "error", err)
			return err
		}
		fmt.Fprintln(fdOut, string(jsonFormattedResult))
	}

	return nil
}
