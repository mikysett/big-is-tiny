package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var flagsTests = []struct {
	description   string
	args          []string
	expectedFlags *Flags
	expectedErr   error
}{
	{
		description:   "Happy path - no flags passed",
		args:          []string{},
		expectedFlags: fixtureFlags(),
	},
	{
		description: "Happy path - all flags passed (long versions)",
		args: []string{
			"-verbose", "-cleanup", "-platform", "azure", "-markdown", "-output", "../file.out", "anotherConfig.json",
		},
		expectedFlags: fixtureFlags(func(f *Flags) {
			f.Cleanup = true
			f.Verbose = true
			f.Platform = Platform(Azure)
			f.MarkdownOut = true
			f.FileOut = "../file.out"
			f.ConfigPath = "anotherConfig.json"
		}),
	},
	{
		description: "Happy path - short versions flags",
		args: []string{
			"-v", "-p", "azure", "-m", "-o", "../file.out",
		},
		expectedFlags: fixtureFlags(func(f *Flags) {
			f.Verbose = true
			f.MarkdownOut = true
			f.FileOut = "../file.out"
			f.Platform = Platform(Azure)
		}),
	},
	{
		description: "Fail on platform flag",
		args: []string{
			"-platform", "invalidPlatform",
		},
		expectedFlags: nil,
		expectedErr:   fmt.Errorf("invalid platform"),
	},
}

func TestGetFlags(t *testing.T) {
	for _, tt := range flagsTests {
		t.Run(tt.description, func(t *testing.T) {
			gotFlags, gotErr := getFlags("bit", tt.args)

			// The Flags result differ
			diff := cmp.Diff(gotFlags, tt.expectedFlags)
			if diff != "" {
				t.Errorf("%v", diff)
			}

			// We get an error when we don't expect it or we don't get one when we expect it
			if tt.expectedErr != nil != (gotErr != nil) {
				t.Errorf("got '%v', want '%v'", gotErr, tt.expectedErr)
			}
		})
	}
}
