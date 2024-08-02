package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var setupConfigTests = []struct {
	description       string
	given             []byte
	expectedBigChange *BigChange
	expectedErr       error
}{
	{
		description:       "Happy path",
		given:             marshalBigChange(fixtureBigChange()),
		expectedBigChange: fixtureBigChange(),
	},
	{
		description: "fail because error unmarshaling",
		given:       []byte("invalid config file"),
		expectedErr: fmt.Errorf("missing or empty config field"),
	},
	{
		description: "fail because mandatory field is missing: Domains",
		given: marshalBigChange(fixtureBigChange(func(bc *BigChange) {
			bc.Domains = nil
		})),
		expectedErr: fmt.Errorf("missing or empty config field"),
	},
	{
		description: "fail because mandatory field is missing: Settings",
		given: marshalBigChange(fixtureBigChange(func(bc *BigChange) {
			bc.Settings = nil
		})),
		expectedErr: fmt.Errorf("missing or empty config field"),
	},
	{
		description: "fail because invalid Settings.MainBranch",
		given: marshalBigChange(fixtureBigChange(func(bc *BigChange) {
			bc.Settings.MainBranch = ""
		})),
		expectedErr: fmt.Errorf("missing or empty config field"),
	},

	{
		description: "fail because invalid Settings.Remote",
		given: marshalBigChange(fixtureBigChange(func(bc *BigChange) {
			bc.Settings.Remote = ""
		})),
		expectedErr: fmt.Errorf("missing or empty config field"),
	},

	{
		description: "fail because invalid Settings.BranchToSplit",
		given: marshalBigChange(fixtureBigChange(func(bc *BigChange) {
			bc.Settings.BranchToSplit = ""
		})),
		expectedErr: fmt.Errorf("missing or empty config field"),
	},
	{
		description: "fail because invalid Domain.Path",
		given: marshalBigChange(fixtureBigChange(func(bc *BigChange) {
			bc.Domains[0].Path = ""
		})),
		expectedErr: fmt.Errorf("missing or empty config field"),
	},
}

func marshalBigChange(bigChange *BigChange) []byte {
	converted, _ := json.Marshal(bigChange)
	return converted
}

func TestSetupConfig(t *testing.T) {
	for _, tt := range setupConfigTests {
		t.Run(tt.description, func(t *testing.T) {
			gotBigChange, gotErr := setupConfig(context.Background(), tt.given)

			// We get an error when we don't expect it or we don't get one when we expect it
			if tt.expectedErr != nil != (gotErr != nil) {
				t.Errorf("got '%v', want '%v'", gotErr, tt.expectedErr)
			}

			// We got a different configuration of what's expected
			diff := cmp.Diff(gotBigChange, tt.expectedBigChange)
			if diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}
