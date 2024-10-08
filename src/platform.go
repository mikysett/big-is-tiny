package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Platform int

const (
	Azure Platform = iota
	GitHub
)

func GetCreatePrForPlatform(p Platform) func(context.Context, *Settings, string, string, string) (string, error) {
	switch p {
	case Platform(GitHub):
		return GitHubCreatePr
	case Platform(Azure):
		return AzureCreatePr
	default:
		panic("unreachable")
	}
}

func GetAbandonPrForPlatform(p Platform) func(context.Context, string) error {
	switch p {
	case Platform(GitHub):
		return GitHubAbandonPr
	case Platform(Azure):
		return AzureAbandonPr
	default:
		panic("unreachable")
	}
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
