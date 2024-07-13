package main

import (
	"encoding/json"
	"fmt"
)

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
