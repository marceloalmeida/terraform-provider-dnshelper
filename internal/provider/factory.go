// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"dnshelper": providerserver.NewProtocol6WithError(New("test")()),
}

var ProtoV6ProviderFactoriesWithEcho = map[string]func() (tfprotov6.ProviderServer, error){
	"dnshelper": providerserver.NewProtocol6WithError(New("test")()),
	"echo":      echoprovider.NewProviderServer(),
}
