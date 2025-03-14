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

func TestDmarcBuilderFunction_Metadata(t *testing.T) {
	f := tffunction.NewDmarcBuilderFunction()
	resp := function.MetadataResponse{}
	f.Metadata(context.Background(), function.MetadataRequest{}, &resp)
	require.Equal(t, "dmarc_builder", resp.Name)
}

func TestDmarcBuilderFunction_Fail_Run(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		wantResp *function.RunResponse
	}{
		{
			name: "Missing inputs",
			args: map[string]interface{}{
				"version":          "DMARC1",
				"policy":           "reject",
				"subdomain_policy": "quarantine",
				"alignment_spf":    "relaxed",
				"alignment_dkim":   "strict",
				"percent":          100,
				"rua":              []string{},
				"ruf":              []string{},
				"failure_options":  "0",
				"report_format":    "afrf",
				//"report_interval":  86400,
			},
			wantErr: true,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.StringValue("")),
				Error:  function.NewFuncError("DmarcBuilder function requires all parameters to be set"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tffunction.NewDmarcBuilderFunction()

			arguments := make([]attr.Value, 10)

			if version, ok := tt.args["version"].(string); ok {
				arguments[0] = types.StringValue(version)
			}
			if policy, ok := tt.args["policy"].(string); ok {
				arguments[1] = types.StringValue(policy)
			}
			if subdomainPolicy, ok := tt.args["subdomain_policy"].(string); ok {
				arguments[2] = types.StringValue(subdomainPolicy)
			}
			if alignmentSPF, ok := tt.args["alignment_spf"].(string); ok {
				arguments[3] = types.StringValue(alignmentSPF)
			}
			if alignmentDKIM, ok := tt.args["alignment_dkim"].(string); ok {
				arguments[4] = types.StringValue(alignmentDKIM)
			}
			if percent, ok := tt.args["percent"].(int); ok {
				arguments[5] = types.Int32Value(int32(percent))
			}
			if rua, ok := tt.args["rua"].([]string); ok {
				arguments[6] = types.ListValueMust(types.StringType, sliceToValues(rua))
			}
			if ruf, ok := tt.args["ruf"].([]string); ok {
				arguments[7] = types.ListValueMust(types.StringType, sliceToValues(ruf))
			}
			if failureOptions, ok := tt.args["failure_options"].(string); ok {
				arguments[8] = types.StringValue(failureOptions)
			}
			if reportFormat, ok := tt.args["report_format"].(string); ok {
				arguments[9] = types.StringValue(reportFormat)
			}
			//if reportInterval, ok := tt.args["report_interval"].(int); ok {
			//	arguments[10] = types.Int32Value(int32(reportInterval))
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

func TestDmarcBuilderFunction_Run(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		wantResp *function.RunResponse
	}{
		{
			name: "valid input",
			args: map[string]interface{}{
				"version":          "DMARC1",
				"policy":           "reject",
				"subdomain_policy": "quarantine",
				"alignment_spf":    "relaxed",
				"alignment_dkim":   "strict",
				"percent":          100,
				"rua":              []string{"mailto:admin@malmeida.dev"},
				"ruf":              []string{"mailto:alerts@malmeida.dev"},
				"failure_options":  "0",
				"report_format":    "afrf",
				"report_interval":  86400,
			},
			wantErr: false,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.StringValue("v=DMARC1; p=reject; sp=quarantine; adkim=s; aspf=r; pct=100; rua=mailto:admin@malmeida.dev; ruf=mailto:alerts@malmeida.dev; fo=0; rf=afrf; ri=86400")),
			},
		},
		{
			name: "invalid policy",
			args: map[string]interface{}{
				"version":          "DMARC1",
				"policy":           "reject",
				"subdomain_policy": "invalid_policy",
				"alignment_spf":    "relaxed",
				"alignment_dkim":   "strict",
				"percent":          100,
				"rua":              []string{"mailto:admin@malmeida.dev"},
				"ruf":              []string{"mailto:alerts@malmeida.dev"},
				"failure_options":  "0",
				"report_format":    "afrf",
				"report_interval":  86400,
			},
			wantErr: true,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.StringValue("")),
				Error:  function.NewFuncError("DmarcBuilder function wrong value for report_interval"),
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			f := tffunction.NewDmarcBuilderFunction()

			arguments := make([]attr.Value, 11)

			if version, ok := tt.args["version"].(string); ok {
				arguments[0] = types.StringValue(version)
			}
			if policy, ok := tt.args["policy"].(string); ok {
				arguments[1] = types.StringValue(policy)
			}
			if subdomainPolicy, ok := tt.args["subdomain_policy"].(string); ok {
				arguments[2] = types.StringValue(subdomainPolicy)
			}
			if alignmentSPF, ok := tt.args["alignment_spf"].(string); ok {
				arguments[3] = types.StringValue(alignmentSPF)
			}
			if alignmentDKIM, ok := tt.args["alignment_dkim"].(string); ok {
				arguments[4] = types.StringValue(alignmentDKIM)
			}
			if percent, ok := tt.args["percent"].(int); ok {
				arguments[5] = types.Int32Value(int32(percent))
			}
			if rua, ok := tt.args["rua"].([]string); ok {
				arguments[6] = types.ListValueMust(types.StringType, sliceToValues(rua))
			}
			if ruf, ok := tt.args["ruf"].([]string); ok {
				arguments[7] = types.ListValueMust(types.StringType, sliceToValues(ruf))
			}
			if failureOptions, ok := tt.args["failure_options"].(string); ok {
				arguments[8] = types.StringValue(failureOptions)
			}
			if reportFormat, ok := tt.args["report_format"].(string); ok {
				arguments[9] = types.StringValue(reportFormat)
			}
			if reportInterval, ok := tt.args["report_interval"].(int); ok {
				arguments[10] = types.Int32Value(int32(reportInterval))
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

func TestAccDmarcBuilderFunction_tf(t *testing.T) {
	t.Setenv("TF_ACC", "1")

	dmarc_version := "DMARC1"
	policy := "reject"
	subdomainPolicy := "quarantine"
	alignmentSPF := "relaxed"
	alignmentDKIM := "strict"
	percent := 100
	rua := []string{"mailto:admin@malmeida.dev"}
	ruf := []string{"mailto:alerts@malmeida.dev"}
	failureOptions := "0"
	reportFormat := "afrf"
	reportInterval := 86400

	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(version.Must(version.NewVersion("1.8.2"))),
			},
			Steps: []resource.TestStep{
				{
					Config: testDmarcBuilderFunctionConfig(dmarc_version, policy, subdomainPolicy, alignmentSPF, alignmentDKIM, percent, rua, ruf, failureOptions, reportFormat, reportInterval),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput(
							"valid_output",
							"v=DMARC1; p=reject; sp=quarantine; adkim=s; aspf=r; pct=100; rua=mailto:admin@malmeida.dev; ruf=mailto:alerts@malmeida.dev; fo=0; rf=afrf; ri=86400",
						),
					),
				},
			},
		},
	)
}

func testDmarcBuilderFunctionConfig(version string, policy string, subdomainPolicy string, alignmentSPF string, alignmentDKIM string, percent int, rua []string, ruf []string, failureOptions string, reportFormat string, reportInterval int) string {
	return fmt.Sprintf(`
		output "valid_output" {
		  value = provider::dnshelper::dmarc_builder(%[1]q, %[2]q, %[3]q, %[4]q, %[5]q, %[6]d, %[7]v, %[8]v, %[9]q, %[10]q, %[11]d)
		}
		`, version, policy, subdomainPolicy, alignmentSPF, alignmentDKIM, percent, types.ListValueMust(types.StringType, sliceToValues(rua)), types.ListValueMust(types.StringType, sliceToValues(ruf)), failureOptions, reportFormat, reportInterval)
}
