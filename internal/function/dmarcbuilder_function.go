// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marceloalmeida/terraform-provider-dnshelper/dnshelper/dmarcbuilder"
)

var (
	_ function.Function = DmarcBuilderFunction{}
)

func NewDmarcBuilderFunction() function.Function {
	return DmarcBuilderFunction{}
}

type DmarcBuilderFunction struct{}

func (r DmarcBuilderFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "dmarc_builder"
}

func (r DmarcBuilderFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "DMARC Builder function",
		MarkdownDescription: "Builds a DMARC record",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "version",
				MarkdownDescription: "The DMARC version, by default DMARC1",
			},
			function.StringParameter{
				Name:                "policy",
				MarkdownDescription: "The DMARC policy (p=), must be one of 'none', 'quarantine', 'reject'",
			},
			function.StringParameter{
				Name:                "subdomain_policy",
				MarkdownDescription: "The DMARC policy for subdomains (sp=), must be one of 'none', 'quarantine', 'reject'",
			},
			function.StringParameter{
				Name:                "alignment_spf",
				MarkdownDescription: "'strict'/'s' or 'relaxed'/'r' alignment for SPF (aspf=, default: 'r')",
			},
			function.StringParameter{
				Name:                "alignment_dkim",
				MarkdownDescription: "'strict'/'s' or 'relaxed'/'r' alignment for DKIM (adkim=, default: 'r')",
			},
			function.Int32Parameter{
				Name:                "percent",
				MarkdownDescription: "Number between 0 and 100, percentage for which policies are applied (pct=, default: 100)",
			},
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "rua",
				MarkdownDescription: "Array of aggregate report targets",
			},
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "ruf",
				MarkdownDescription: "Array of failure report targets",
			},
			function.StringParameter{
				Name:                "failure_options",
				MarkdownDescription: "String containing is passed raw (fo=, default: '0')",
			},
			function.StringParameter{
				Name:                "failure_format",
				MarkdownDescription: "Format in which failure reports are requested (rf=, default: 'afrf')",
			},
			function.Int32Parameter{
				Name:                "report_interval",
				MarkdownDescription: "Interval in which reports are requested (ri=)",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r DmarcBuilderFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data struct {
		Version         string   `tfsdk:"version"`
		Policy          string   `tfsdk:"policy"`
		SubdomainPolicy string   `tfsdk:"subdomain_policy"`
		AlignmentSPF    string   `tfsdk:"alignment_spf"`
		AlignmentDKIM   string   `tfsdk:"alignment_dkim"`
		Percent         int32    `tfsdk:"percent"`
		RUA             []string `tfsdk:"rua"`
		RUF             []string `tfsdk:"ruf"`
		FailureOptions  string   `tfsdk:"failure_options"`
		FailureFormat   string   `tfsdk:"failure_format"`
		ReportInterval  int32    `tfsdk:"report_interval"`
	}

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &data.Version, &data.Policy, &data.SubdomainPolicy, &data.AlignmentSPF, &data.AlignmentDKIM, &data.Percent, &data.RUA, &data.RUF, &data.FailureOptions, &data.FailureFormat, &data.ReportInterval))

	if resp.Error != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Errorf("failed to get arguments").Error()))
		return
	}

	config := dmarcbuilder.DMARCConfig{
		Version:         data.Version,
		Policy:          data.Policy,
		SubdomainPolicy: data.SubdomainPolicy,
		AlignmentSPF:    data.AlignmentSPF,
		AlignmentDKIM:   data.AlignmentDKIM,
		Percent:         data.Percent,
		RUA:             data.RUA,
		RUF:             data.RUF,
		FailureOptions:  data.FailureOptions,
		FailureFormat:   data.FailureFormat,
		ReportInterval:  data.ReportInterval,
	}
	result, err := dmarcbuilder.DmarcBuilder(config)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, &result))
}
