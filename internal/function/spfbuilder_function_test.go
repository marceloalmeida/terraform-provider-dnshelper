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
	"github.com/stretchr/testify/require"

	tffunction "github.com/marceloalmeida/terraform-provider-dnshelper/internal/function"
	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/provider"
)

func TestSPFBuilderFunction_Metadata(t *testing.T) {
	f := tffunction.NewSPFBuilderFunction()
	resp := function.MetadataResponse{}
	f.Metadata(context.Background(), function.MetadataRequest{}, &resp)
	require.Equal(t, "spf_builder", resp.Name)
}

func TestSPFBuilderFunction_Definition(t *testing.T) {
	f := tffunction.NewSPFBuilderFunction()
	resp := function.DefinitionResponse{}
	f.Definition(context.Background(), function.DefinitionRequest{}, &resp)
	require.Equal(t, "SPF Builder function", resp.Definition.Summary)
	require.Equal(t, "Builds an SPF record", resp.Definition.MarkdownDescription)
	require.Len(t, resp.Definition.Parameters, 6)
	require.Equal(t, "domain", resp.Definition.Parameters[0].GetName())
	require.Equal(t, "overflow", resp.Definition.Parameters[1].GetName())
	require.Equal(t, "txt_max_size", resp.Definition.Parameters[2].GetName())
	require.Equal(t, "domain_on_record_key", resp.Definition.Parameters[3].GetName())
	require.Equal(t, "parts", resp.Definition.Parameters[4].GetName())
	require.Equal(t, "flatten", resp.Definition.Parameters[5].GetName())
	require.Equal(t, types.ListType{ElemType: types.StringType}, resp.Definition.Parameters[4].GetType())
}

func TestSPFBuilderFunction_Run(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		wantErr  bool
		wantResp *function.RunResponse
	}{
		{
			name: "valid input",
			args: map[string]interface{}{
				"domain":               "example.com",
				"overflow":             "spf%d",
				"txt_max_size":         255,
				"domain_on_record_key": true,
				"parts":                []string{"v=spf1", "include:_spf.google.com", "~all"},
				"flatten":              []string{"example.com"},
			},
			wantErr: false,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.MapValueMust(
					types.ListType{ElemType: types.StringType},
					map[string]attr.Value{"example.com": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("v=spf1 include:_spf.google.com ~all")})},
				)),
				Error: nil,
			},
		},
		{
			name: "invalid txtMaxSize",
			args: map[string]interface{}{
				"domain":               "example.com",
				"overflow":             "spf%d",
				"txt_max_size":         0,
				"domain_on_record_key": true,
				"parts":                []string{"v=spf1", "include:_spf.google.com", "~all"},
				"flatten":              []string{"example.com"},
			},
			wantErr: true,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.ObjectNull(map[string]attr.Type{
					"records": types.ListType{ElemType: types.StringType},
				})),
			},
		},
		{
			name: "invalid overflow format",
			args: map[string]interface{}{
				"domain":               "example.com",
				"overflow":             "spf",
				"txt_max_size":         255,
				"domain_on_record_key": true,
				"parts":                []string{"v=spf1", "include:_spf.google.com", "~all"},
				"flatten":              []string{"example.com"},
			},
			wantErr: true,
			wantResp: &function.RunResponse{
				Result: function.NewResultData(types.ObjectNull(map[string]attr.Type{
					"records": types.ListType{ElemType: types.StringType},
				})),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tffunction.NewSPFBuilderFunction()

			domain, ok := tt.args["domain"].(string)
			if !ok {
				t.Fatal("domain is not a string")
			}
			overflow, ok := tt.args["overflow"].(string)
			if !ok {
				t.Fatal("overflow is not a string")
			}
			txtMaxSize, ok := tt.args["txt_max_size"].(int)
			if !ok {
				t.Fatal("txtMaxSize is not an int")
			}
			domainOnRecordKey, ok := tt.args["domain_on_record_key"].(bool)
			if !ok {
				t.Fatal("domain_on_record_key is not a bool")
			}
			parts, ok := tt.args["parts"].([]string)
			if !ok {
				t.Fatal("parts is not a []string")
			}
			flatten, ok := tt.args["flatten"].([]string)
			if !ok {
				t.Fatal("flatten is not a []string")
			}

			req := function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(domain),
					types.StringValue(overflow),
					types.Int32Value(int32(txtMaxSize)),
					types.BoolValue(domainOnRecordKey),
					types.ListValueMust(types.StringType, sliceToValues(parts)),
					types.ListValueMust(types.StringType, sliceToValues(flatten)),
				}),
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

			require.NotNil(t, resp.Result)
		})
	}
}

func TestAccSPFBuilderFunction_tf(t *testing.T) {
	t.Setenv("TF_ACC", "1")

	domain := "example.com"
	overflow := "spf%d"
	txtMaxSize := 255
	domainOnRecordKey := true
	parts := []string{"v=spf1", "include:_spf.google.com", "~all"}
	flatten := []string{"example.com"}

	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: provider.ProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(version.Must(version.NewVersion("1.8.2"))),
			},
			Steps: []resource.TestStep{
				{
					Config: testSpfBuilderFunctionConfig(domain, overflow, txtMaxSize, domainOnRecordKey, parts, flatten),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput(
							"valid_output_jsonencode",
							`{"@":["v=spf1 include:_spf.google.com ~all"]}`,
						),
					),
				},
			},
		},
	)
}
func TestSPFBuilderFunction_Run_Error(t *testing.T) {
	testCases := []struct {
		name        string
		args        []attr.Value
		expectError bool
	}{
		{
			name: "incorrect argument types",
			args: []attr.Value{
				types.Int64Value(123), // domain should be string
				types.StringValue("spf%d"),
				types.Int32Value(255),
				types.BoolValue(true),
				types.ListValueMust(types.StringType, []attr.Value{}),
				types.ListValueMust(types.StringType, []attr.Value{}),
			},
			expectError: true,
		},
		{
			name: "missing arguments",
			args: []attr.Value{
				types.StringValue("example.com"),
				types.StringValue("spf%d"),
			},
			expectError: true,
		},
		{
			name: "invalid argument count",
			args: []attr.Value{
				types.StringValue("example.com"),
				types.StringValue("spf%d"),
				types.Int32Value(255),
				types.BoolValue(true),
				types.ListValueMust(types.StringType, []attr.Value{}),
				types.ListValueMust(types.StringType, []attr.Value{}),
				types.StringValue("extra"), // Extra argument
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := tffunction.NewSPFBuilderFunction()
			resp := &function.RunResponse{}

			f.Run(context.Background(), function.RunRequest{
				Arguments: function.NewArgumentsData(tc.args),
			}, resp)

			if tc.expectError && resp.Error == nil {
				t.Fatal("expected error but got none")
			}
			if !tc.expectError && resp.Error != nil {
				t.Fatalf("expected no error but got: %v", resp.Error)
			}
		})
	}
}
func testSpfBuilderFunctionConfig(domain string, overflow string, txtMaxSize int, domainOnRecordKey bool, parts []string, flatten []string) string {

	return fmt.Sprintf(`
output "valid_output_jsonencode" {
  value = jsonencode(provider::dnshelper::spf_builder(%[1]q, %[2]q, %[3]d, %[4]t, %[5]v, %[6]v))
}
`, domain, overflow, txtMaxSize, domainOnRecordKey, types.ListValueMust(types.StringType, sliceToValues(parts)), sliceToValues(flatten))
}

func sliceToValues(slice []string) []attr.Value {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return values
}
