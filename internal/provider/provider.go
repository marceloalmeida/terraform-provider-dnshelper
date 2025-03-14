// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	tffunction "github.com/marceloalmeida/terraform-provider-dnshelper/internal/function"
)

var _ provider.Provider = &DnshelperProvider{}
var _ provider.ProviderWithFunctions = &DnshelperProvider{}
var _ provider.ProviderWithEphemeralResources = &DnshelperProvider{}

type DnshelperProvider struct {
	version string
}

type DnshelperProviderModel struct{}

func (p *DnshelperProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dnshelper"
	resp.Version = p.version
}

func (p *DnshelperProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *DnshelperProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DnshelperProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DnshelperProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *DnshelperProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *DnshelperProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *DnshelperProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		tffunction.NewSPFBuilderFunction,
		tffunction.NewCAABuilderFunction,
		tffunction.NewDmarcBuilderFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DnshelperProvider{
			version: version,
		}
	}
}
