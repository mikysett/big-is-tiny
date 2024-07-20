package main

import (
	"context"
	"encoding/json"
	"fmt"
	"syscall"
)

func chdirWithLogs(ctx context.Context, path string) error {
	err := syscall.Chdir(path)
	if err != nil {
		log := LoggerFromContext(ctx)
		log.Error("failed to change directory", "target directory", path, "error", err)
		return err
	}
	return nil
}

func (e Communication) String() string {
	switch e {
	case Slack:
		return "Slack"
	case Teams:
		return "Teams"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func (s Communication) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (e Platform) String() string {
	switch e {
	case Azure:
		return "Azure"
	case GitHub:
		return "GitHub"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func (s Platform) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
