package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type givenGenerateFromTemplate struct {
	domain   *Domain
	template string
}

var generateFromTemplateTests = []struct {
	description    string
	given          *givenGenerateFromTemplate
	expectedResult string
}{
	{
		description: "Happy path",
		given: &givenGenerateFromTemplate{
			domain: &Domain{
				Id:   "AA",
				Name: "backend",
				Teams: []Team{
					{
						Name: "principal",
						Url:  "team1.com",
					},
					{
						Name: "secondary",
						Url:  "team2.com",
					},
				},
			},
			template: "[{{change_id}}] {{domain_id}} {{domain_name}}: Big change split {{domain_name}}, {{team_name_1}}({{team_url_1}}), {{team_name_2}}({{team_url_2}})\n{{team_name_3}}({{team_url_3}})",
		},
		expectedResult: "[BIT001] AA backend: Big change split backend, principal(team1.com), secondary(team2.com)\n{{team_name_3}}({{team_url_3}})",
	},
}

func TestGenerateFromTemplate(t *testing.T) {
	for _, tt := range generateFromTemplateTests {
		svc := &BigChange{Id: "BIT001"}

		gotResult := svc.generateFromTemplate(tt.given.domain, tt.given.template)

		// We got a different result of what's expected
		diff := cmp.Diff(gotResult, tt.expectedResult)
		if diff != "" {
			t.Errorf("%v", diff)
		}
	}
}
