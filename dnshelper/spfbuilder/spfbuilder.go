// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package spfbuilder

import (
	"fmt"
	"strings"

	"github.com/StackExchange/dnscontrol/v4/pkg/spflib"
)

type Resolver interface {
	GetTXT(domain string) ([]string, error)
}

func BuildSPFRecord(domain string, overflow string, txtMaxSize int32, domainOnRecordKey bool, parts []string, flatten []string) (map[string][]string, error) {
	return BuildSPFRecordWithResolver(domain, overflow, txtMaxSize, domainOnRecordKey, parts, flatten, &spflib.LiveResolver{})
}

func BuildSPFRecordWithResolver(domain string, overflow string, txtMaxSize int32, domainOnRecordKey bool, parts []string, flatten []string, resolver spflib.Resolver) (map[string][]string, error) {
	spfRecord := strings.Join(parts, " ")
	rec, err := spflib.Parse(spfRecord, resolver)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SPF record: %w", err)
	}

	if txtMaxSize < 1 {
		return nil, fmt.Errorf("txtMaxSize must be greater than 0")
	}

	if !strings.Contains(overflow, "%d") {
		return nil, fmt.Errorf("split format `%s` in `%s` is not proper format (missing `%%d`)", overflow, domain)
	}

	for _, domain := range flatten {
		rec = rec.Flatten(domain)
	}

	rec = dedup(rec)

	splitRec := rec.TXTSplit(overflow+"."+domain, 0, int(txtMaxSize))

	result := make(map[string][]string, len(splitRec))
	for k, v := range splitRec {
		if !domainOnRecordKey {
			k = strings.TrimSuffix(k, "."+domain)
		}
		result[k] = v
	}

	return result, nil
}

func dedup(s *spflib.SPFRecord) *spflib.SPFRecord {
	seen := map[string]bool{}
	newParts := make([]*spflib.SPFPart, 0, len(s.Parts))
	for _, p := range s.Parts {
		if seen[p.Text] {
			continue
		}
		seen[p.Text] = true
		newParts = append(newParts, p)
	}
	s.Parts = newParts
	return s
}
