// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/caabuilder"
)

var (
	_ function.Function = CAABuilderFunction{}
)

func NewCAABuilderFunction() function.Function {
	return CAABuilderFunction{}
}

type CAABuilderFunction struct{}

func (r CAABuilderFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "caa_builder"
}

func (r CAABuilderFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "CAA Builder function",
		MarkdownDescription: "Builds a CAA records",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "iodef",
				MarkdownDescription: "The contact mail address",
			},
			function.BoolParameter{
				Name:                "iodef_critical",
				MarkdownDescription: "Boolean if sending report is required/critical",
			},
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "issue",
				MarkdownDescription: "List of CAs which are allowed to issue certificates for the domain",
			},
			function.BoolParameter{
				Name:                "issue_critical",
				MarkdownDescription: "Boolean if issue is required/critical",
			},
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "issuewild",
				MarkdownDescription: "Allowed CAs which can issue wildcard certificates for this domain",
			},
			function.BoolParameter{
				Name:                "issuewild_critical",
				MarkdownDescription: "Boolean if issuewild is required/critical",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (r CAABuilderFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data struct {
		Iodef             string   `tfsdk:"iodef"`
		IodefCritical     bool     `tfsdk:"iodef_critical"`
		Issue             []string `tfsdk:"issue"`
		IssueCritical     bool     `tfsdk:"issue_critical"`
		Issuewild         []string `tfsdk:"issuewild"`
		IssuewildCritical bool     `tfsdk:"issuewild_critical"`
	}

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &data.Iodef, &data.IodefCritical, &data.Issue, &data.IssueCritical, &data.Issuewild, &data.IssuewildCritical))

	if resp.Error != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Errorf("failed to get arguments").Error()))
		return
	}

	config := caabuilder.CAAConfig{
		Iodef:             data.Iodef,
		IodefCritical:     data.IodefCritical,
		Issue:             data.Issue,
		IssueCritical:     data.IssueCritical,
		Issuewild:         data.Issuewild,
		IssuewildCritical: data.IssuewildCritical,
	}
	result, err := caabuilder.CAABuilderString(config)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, &result))
}
