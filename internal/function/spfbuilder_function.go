// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/spfbuilder"
	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/testutil"
)

var (
	_ function.Function = SPFBuilderFunction{}
)

func NewSPFBuilderFunction() function.Function {
	return SPFBuilderFunction{}
}

type SPFBuilderFunction struct{}

func (r SPFBuilderFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "spf_builder"
}

func (r SPFBuilderFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "SPF Builder function",
		MarkdownDescription: "Builds an SPF record",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "domain",
				MarkdownDescription: "Domain to build the SPF record for",
			},
			function.StringParameter{
				Name:                "overflow",
				MarkdownDescription: "Overflow value",
			},
			function.Int32Parameter{
				Name:                "txt_max_size",
				MarkdownDescription: "TXT max size",
			},
			function.BoolParameter{
				Name:                "domain_on_record_key",
				MarkdownDescription: "Whether to include the TLD on the record key",
			},
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "parts",
				MarkdownDescription: "SPF parts",
			},
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "flatten",
				MarkdownDescription: "A list of domains to flatten",
			},
		},
		Return: function.MapReturn{
			ElementType: types.ListType{ElemType: types.StringType},
		},
	}
}

func (r SPFBuilderFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data struct {
		Domain            string   `tfsdk:"domain"`
		Overflow          string   `tfsdk:"overflow"`
		TxtMaxSize        int32    `tfsdk:"txt_max_size"`
		DomainOnRecordKey bool     `tfsdk:"domain_on_record_key"`
		Parts             []string `tfsdk:"parts"`
		Flatten           []string `tfsdk:"flatten"`
	}

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &data.Domain, &data.Overflow, &data.TxtMaxSize, &data.DomainOnRecordKey, &data.Parts, &data.Flatten))

	if resp.Error != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Errorf("failed to get arguments").Error()))
		return
	}

	result, err := buildSPFRecord(data.Domain, data.Overflow, data.TxtMaxSize, data.DomainOnRecordKey, data.Parts, data.Flatten)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, &result))
}

func buildSPFRecord(domain string, overflow string, txtMaxSize int32, domainOnRecordKey bool, parts []string, flatten []string) (map[string][]string, error) {
	if testing.Testing() && !strings.HasPrefix(os.Getenv("TF_ACC"), "1") {
		mock := testutil.NewMockResolver()
		return spfbuilder.BuildSPFRecordWithResolver(domain, overflow, txtMaxSize, domainOnRecordKey, parts, flatten, mock)
	}
	return spfbuilder.BuildSPFRecord(domain, overflow, txtMaxSize, domainOnRecordKey, parts, flatten)
}
