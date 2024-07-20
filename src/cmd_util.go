package main

import (
	"context"
	"os/exec"
	"strings"
)

func runCmd(ctx context.Context, name string, args ...string) ([]byte, error) {
	log := LoggerFromContext(ctx)

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("failed to run command",
			"command", name,
			"args", strings.Join(args, " "),
			"output", string(output[:]),
			"error", err)
	} else {
		log.Debug("run command",
			"command", name,
			"args", strings.Join(args, " "),
			"output", string(output[:]),
			"error", err)
	}
	return output, err
}
