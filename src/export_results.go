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

	if config.Settings.OutputTemplate != "" {
		for _, domain := range config.Domains {
			fmt.Fprintf(fdOut, "%s\n", config.generateFromTemplate(domain, config.Settings.OutputTemplate))
		}
	} else {
		// Default writes json formatted Branch name and PR URL
		createdPrs := make([]createdPr, 0, len(config.Domains))
		for _, domain := range config.Domains {
			if domain.PullRequest.Url == "" {
				continue
			}
			createdPrs = append(createdPrs, createdPr{
				Branch: domain.Branch.Name,
				PrUrl:  domain.PullRequest.Url,
			})
		}

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
