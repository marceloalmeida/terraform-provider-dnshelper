// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"context"
	"testing"

	tpf "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/provider"
)

func TestProvider(t *testing.T) {
	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}
}

func TestProviderSchema(t *testing.T) {
	t.Parallel()

	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}
}

func TestProviderResources(t *testing.T) {
	t.Parallel()

	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}

	resources := p.Resources(context.Background())
	if resources == nil {
		t.Fatal("resources should not be nil")
	}
}

func TestProviderDataSources(t *testing.T) {
	t.Parallel()

	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}

	dataSources := p.DataSources(context.Background())
	if dataSources == nil {
		t.Fatal("data sources should not be nil")
	}
}
func TestProviderMetadata(t *testing.T) {
	t.Parallel()

	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}

	var metadataResp tpf.MetadataResponse
	p.Metadata(context.Background(), tpf.MetadataRequest{}, &metadataResp)

	if metadataResp.TypeName != "dnshelper" {
		t.Errorf("expected provider type name to be 'dnshelper', got %s", metadataResp.TypeName)
	}

	if metadataResp.Version != "test" {
		t.Errorf("expected provider version to be 'test', got %s", metadataResp.Version)
	}
}
func TestProviderSchemaEmpty(t *testing.T) {
	t.Parallel()

	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}

	var schemaResp tpf.SchemaResponse
	p.Schema(context.Background(), tpf.SchemaRequest{}, &schemaResp)

	if len(schemaResp.Schema.Attributes) != 0 {
		t.Errorf("expected provider schema to have 0 attributes, got %d", len(schemaResp.Schema.Attributes))
	}

	if len(schemaResp.Schema.Blocks) != 0 {
		t.Errorf("expected provider schema to have 0 blocks, got %d", len(schemaResp.Schema.Blocks))
	}
}
func TestProviderFunctions(t *testing.T) {
	t.Parallel()

	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}

	providerWithFunctions, ok := p.(tpf.ProviderWithFunctions)
	if !ok {
		t.Fatal("provider does not implement ProviderWithFunctions")
	}
	functions := providerWithFunctions.Functions(context.Background())
	if functions == nil {
		t.Fatal("functions should not be nil")
	}

	for i, fn := range functions {
		f := fn()
		if f == nil {
			t.Fatalf("function %d should not be nil", i)
		}
	}
}

func TestProviderEphemeralResources(t *testing.T) {
	t.Parallel()

	p := provider.New("test")()
	if p == nil {
		t.Fatal("failed to create provider")
	}

	providerWithEphemeralResources, ok := p.(tpf.ProviderWithEphemeralResources)
	if !ok {
		t.Fatal("provider does not implement ProviderWithEphemeralResources")
	}
	ephemeralResources := providerWithEphemeralResources.EphemeralResources(context.Background())
	if ephemeralResources == nil {
		t.Fatal("Ephemeral Resources should not be nil")
	}

	for i, er := range ephemeralResources {
		f := er()
		if f == nil {
			t.Fatalf("Ephemeral Resource %d should not be nil", i)
		}
	}
}
