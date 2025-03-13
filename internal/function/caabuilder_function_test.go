// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package function_test

import (
	"context"
	"fmt"
	"testing"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/provider"
	"github.com/stretchr/testify/require"

	tffunction "github.com/marceloalmeida/terraform-provider-dnshelper/internal/function"
)

func TestCaaBuilderFunction_Metadata(t *testing.T) {
	f := tffunction.NewCAABuilderFunction()
	resp := function.MetadataResponse{}
	f.Metadata(context.Background(), function.MetadataRequest{}, &resp)
	require.Equal(t, "caa_builder", resp.Name)
}

func TestCaaBuilderFunction_Definition(t *testing.T) {
	f := tffunction.NewCAABuilderFunction()
	resp := function.DefinitionResponse{}
	f.Definition(context.Background(), function.DefinitionRequest{}, &resp)
	require.Equal(t, "CAA Builder function", resp.Definition.Summary)
	require.Equal(t, "Builds a CAA records", resp.Definition.MarkdownDescription)
	require.Len(t, resp.Definition.Parameters, 6)
	require.Equal(t, "iodef", resp.Definition.Parameters[0].GetName())
	require.Equal(t, "iodef_critical", resp.Definition.Parameters[1].GetName())
	require.Equal(t, "issue", resp.Definition.Parameters[2].GetName())
	require.Equal(t, "issue_critical", resp.Definition.Parameters[3].GetName())
	require.Equal(t, "issuewild", resp.Definition.Parameters[4].GetName())
	require.Equal(t, "issuewild_critical", resp.Definition.Parameters[5].GetName())
	require.Equal(t, types.ListType{ElemType: types.StringType}, resp.Definition.Parameters[2].GetType())
	require.Equal(t, types.ListType{ElemType: types.StringType}, resp.Definition.Parameters[4].GetType())
}

func TestCaaBuilderFunction_Fail_Run(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		wantResp *function.RunResponse
	}{
		{
			name: "Missing inputs",
			args: map[string]interface{}{
				"iodef":          "asd",
				"iodef_critical": true,
				"issue":          []string{},
				"issue_critical": false,
				"issuewild":      []string{},
				//"issuewild_critical": false,
			},
			wantErr: true,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.ListValueMust(types.StringType, []attr.Value{})),
				Error:  function.NewFuncError("CAABuilder requires at least one entry in issue or issuewild"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tffunction.NewCAABuilderFunction()

			arguments := make([]attr.Value, 5)

			if iodef, ok := tt.args["iodef"].(string); ok {
				arguments[0] = types.StringValue(iodef)
			}
			if iodefCritical, ok := tt.args["iodef_critical"].(bool); ok {
				arguments[1] = types.BoolValue(iodefCritical)
			}
			if issue, ok := tt.args["issue"].([]string); ok {
				arguments[2] = types.ListValueMust(types.StringType, sliceToValues(issue))
			}
			if issueCritical, ok := tt.args["issue_critical"].(bool); ok {
				arguments[3] = types.BoolValue(issueCritical)
			}
			if issuewild, ok := tt.args["issuewild"].([]string); ok {
				arguments[4] = types.ListValueMust(types.StringType, sliceToValues(issuewild))
			}
			//if issuewildCritical, ok := tt.args["issuewild_critical"].(bool); ok {
			//	arguments[5] = types.BoolValue(issuewildCritical)
			//}

			req := function.RunRequest{
				Arguments: function.NewArgumentsData(arguments),
			}

			resp := tt.wantResp
			f.Run(context.Background(), req, resp)

			if tt.wantErr {
				require.Error(t, resp.Error)
				return
			}

			if resp.Error.Equal(nil) {
				require.NoError(t, nil)
			}

			require.Equal(t, tt.wantResp, resp)
		})
	}
}

func TestCaaBuilderFunction_Run(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		wantResp *function.RunResponse
	}{
		{
			name: "valid input",
			args: map[string]interface{}{
				"iodef":              "mailto:domain-names@malmeida.dev",
				"iodef_critical":     true,
				"issue":              []string{"amazon.com", "comodoca.com", "digicert.com; cansignhttpexchanges=yes", "letsencrypt.org", "pki.goog; cansignhttpexchanges=yes", "sectigo.com", "ssl.com"},
				"issue_critical":     false,
				"issuewild":          []string{"amazon.com", "comodoca.com", "digicert.com; cansignhttpexchanges=yes", "letsencrypt.org", "pki.goog; cansignhttpexchanges=yes", "sectigo.com", "ssl.com"},
				"issuewild_critical": false,
			},
			wantErr: false,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("128 iodef \"mailto:domain-names@malmeida.dev\""),
					types.StringValue("0 issue \"amazon.com\""),
					types.StringValue("0 issue \"comodoca.com\""),
					types.StringValue("0 issue \"digicert.com; cansignhttpexchanges=yes\""),
					types.StringValue("0 issue \"letsencrypt.org\""),
					types.StringValue("0 issue \"pki.goog; cansignhttpexchanges=yes\""),
					types.StringValue("0 issue \"sectigo.com\""),
					types.StringValue("0 issue \"ssl.com\""),
					types.StringValue("0 issuewild \"amazon.com\""),
					types.StringValue("0 issuewild \"comodoca.com\""),
					types.StringValue("0 issuewild \"digicert.com; cansignhttpexchanges=yes\"")},
				)),
				Error: nil,
			},
		},
		{
			name: "invalid input",
			args: map[string]interface{}{
				"iodef":              "asd",
				"iodef_critical":     true,
				"issue":              []string{},
				"issue_critical":     false,
				"issuewild":          []string{},
				"issuewild_critical": false,
			},
			wantErr: true,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.ListValueMust(types.StringType, []attr.Value{})),
				Error:  function.NewFuncError("CAABuilder requires at least one entry in issue or issuewild"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tffunction.NewCAABuilderFunction()

			arguments := make([]attr.Value, 6)

			if iodef, ok := tt.args["iodef"].(string); ok {
				arguments[0] = types.StringValue(iodef)
			}
			if iodefCritical, ok := tt.args["iodef_critical"].(bool); ok {
				arguments[1] = types.BoolValue(iodefCritical)
			}
			if issue, ok := tt.args["issue"].([]string); ok {
				arguments[2] = types.ListValueMust(types.StringType, sliceToValues(issue))
			}
			if issueCritical, ok := tt.args["issue_critical"].(bool); ok {
				arguments[3] = types.BoolValue(issueCritical)
			}
			if issuewild, ok := tt.args["issuewild"].([]string); ok {
				arguments[4] = types.ListValueMust(types.StringType, sliceToValues(issuewild))
			}
			if issuewildCritical, ok := tt.args["issuewild_critical"].(bool); ok {
				arguments[5] = types.BoolValue(issuewildCritical)
			}

			req := function.RunRequest{
				Arguments: function.NewArgumentsData(arguments),
			}

			resp := tt.wantResp
			f.Run(context.Background(), req, resp)

			if tt.wantErr {
				require.Error(t, resp.Error)
				return
			}

			if resp.Error.Equal(nil) {
				require.NoError(t, nil)
			}

			require.Equal(t, tt.wantResp, resp)
		})
	}
}

func TestAccCaaBuilderFunction_tf(t *testing.T) {
	t.Setenv("TF_ACC", "1")

	iodef := "mailto:domain-names@malmeida.dev"
	iodefCritical := true
	issue := []string{"amazon.com", "comodoca.com", "digicert.com; cansignhttpexchanges=yes", "letsencrypt.org", "pki.goog; cansignhttpexchanges=yes", "sectigo.com", "ssl.com"}
	issueCritical := false
	issuewild := []string{"amazon.com", "comodoca.com", "digicert.com; cansignhttpexchanges=yes", "letsencrypt.org", "pki.goog; cansignhttpexchanges=yes", "sectigo.com", "ssl.com"}
	issuewildCritical := false

	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(version.Must(version.NewVersion("1.8.2"))),
			},
			Steps: []resource.TestStep{
				{
					Config: testCaaBuilderFunctionConfig(iodef, iodefCritical, issue, issueCritical, issuewild, issuewildCritical),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput(
							"valid_output_jsonencode",
							"[\"128 iodef \\\"mailto:domain-names@malmeida.dev\\\"\",\"0 issue \\\"amazon.com\\\"\",\"0 issue \\\"comodoca.com\\\"\",\"0 issue \\\"digicert.com; cansignhttpexchanges=yes\\\"\",\"0 issue \\\"letsencrypt.org\\\"\",\"0 issue \\\"pki.goog; cansignhttpexchanges=yes\\\"\",\"0 issue \\\"sectigo.com\\\"\",\"0 issue \\\"ssl.com\\\"\",\"0 issuewild \\\"amazon.com\\\"\",\"0 issuewild \\\"comodoca.com\\\"\",\"0 issuewild \\\"digicert.com; cansignhttpexchanges=yes\\\"\",\"0 issuewild \\\"letsencrypt.org\\\"\",\"0 issuewild \\\"pki.goog; cansignhttpexchanges=yes\\\"\",\"0 issuewild \\\"sectigo.com\\\"\",\"0 issuewild \\\"ssl.com\\\"\"]",
						),
					),
				},
			},
		},
	)
}

func testCaaBuilderFunctionConfig(iodef string, iodefCritical bool, issue []string, issueCritical bool, issuewild []string, issuewildCritical bool) string {
	return fmt.Sprintf(`
output "valid_output_jsonencode" {
  value = jsonencode(provider::dnshelper::caa_builder(%[1]q, %[2]t, %[3]v, %[4]t, %[5]v, %[6]t))
}
`, iodef, iodefCritical, types.ListValueMust(types.StringType, sliceToValues(issue)), issueCritical, types.ListValueMust(types.StringType, sliceToValues(issuewild)), issuewildCritical)
}
