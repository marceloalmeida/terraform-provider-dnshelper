// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"log"

	"github.com/StackExchange/dnscontrol/v4/pkg/spflib"
)

type MockResolver struct {
	TxtRecords map[string][]string
}

func (m *MockResolver) GetTXT(domain string) ([]string, error) {
	if records, ok := m.TxtRecords[domain]; ok {
		return records, nil
	}
	return nil, nil
}

func (m *MockResolver) GetSPF(domain string) (string, error) {
	if records, ok := m.TxtRecords[domain]; ok && len(records) > 0 {
		return records[0], nil
	}
	return "", nil
}

func NewMockResolver() spflib.Resolver {
	res, err := spflib.NewCache("../../internal/testutil/testdata-dns.json")
	if err != nil {
		log.Fatalf("error creating mock resolver: %v", err)
		return nil
	}
	return res
}
