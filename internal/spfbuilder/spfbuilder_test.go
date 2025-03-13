// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package spfbuilder_test

import (
	"testing"

	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/spfbuilder"
	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/testutil"
)

func TestBuildSPFRecord(t *testing.T) {
	mock := testutil.NewMockResolver()

	tests := []struct {
		name              string
		domain            string
		overflow          string
		txtMaxSize        int32
		domainOnRecordKey bool
		parts             []string
		flatten           []string
		want              map[string][]string
		wantErr           bool
	}{
		{
			name:              "basic SPF record",
			domain:            "example.com",
			overflow:          "spf%d",
			txtMaxSize:        255,
			domainOnRecordKey: false,
			parts:             []string{"v=spf1", "ip4:192.0.2.0/24", "include:_spf.example.com", "-all"},
			flatten:           []string{},
			want: map[string][]string{
				"@": {"v=spf1 ip4:192.0.2.0/24 include:_spf.example.com -all"},
			},
			wantErr: false,
		},
		{
			name:              "SPF record with flattening",
			domain:            "example.com",
			overflow:          "spf%d",
			txtMaxSize:        255,
			domainOnRecordKey: false,
			parts:             []string{"v=spf1", "include:example.com", "-all"},
			flatten:           []string{"example.com", "_spf.example.com"},
			want: map[string][]string{
				"@": {"v=spf1 ip4:192.168.2.1/32 -all"},
			},
			wantErr: false,
		},
		{
			name:              "SPF record with domain on record key",
			domain:            "example.org",
			overflow:          "spf%d",
			txtMaxSize:        100,
			domainOnRecordKey: true,
			parts:             []string{"v=spf1", "include:example.org", "~all"},
			flatten:           []string{"example.org", "_spf.example.org"},
			want: map[string][]string{
				"@":                {"v=spf1 ip4:192.168.0.1/32 ip4:192.168.0.2/32 ip4:192.168.0.3/32 include:spf1.example.org ~all"},
				"spf1.example.org": {"v=spf1 ip4:192.168.0.4/32 ip4:192.168.0.5/32 ip4:192.168.0.6/32 include:spf2.example.org ~all"},
				"spf2.example.org": {"v=spf1 ip4:192.168.0.7/32 ip4:192.168.0.8/32 ip4:192.168.0.9/32 include:spf3.example.org ~all"},
				"spf3.example.org": {"v=spf1 ip4:192.168.0.10/32 ip4:192.168.0.100/32 ip4:192.168.0.200/32 ip6:fe80:831e:c000::/38 ~all"},
			},
			wantErr: false,
		},
		{
			name:              "SPF record with splitting",
			domain:            "example.org",
			overflow:          "spf%d",
			txtMaxSize:        255,
			domainOnRecordKey: false,
			parts:             []string{"v=spf1", "include:example.org", "~all"},
			flatten:           []string{"example.org", "_spf.example.org"},
			want: map[string][]string{
				"@":    {"v=spf1 ip4:192.168.0.1/32 ip4:192.168.0.2/32 ip4:192.168.0.3/32 ip4:192.168.0.4/32 ip4:192.168.0.5/32 ip4:192.168.0.6/32 ip4:192.168.0.7/32 ip4:192.168.0.8/32 ip4:192.168.0.9/32 ip4:192.168.0.10/32 ip4:192.168.0.100/32 include:spf1.example.org ~all"},
				"spf1": {"v=spf1 ip4:192.168.0.200/32 ip6:fe80:831e:c000::/38 ~all"},
			},
			wantErr: false,
		},
		{
			name:              "invalid txtMaxSize",
			domain:            "example.com",
			overflow:          "spf%d",
			txtMaxSize:        0,
			domainOnRecordKey: false,
			parts:             []string{"v=spf1", "ip4:192.0.2.0/24", "-all"},
			flatten:           []string{},
			want:              nil,
			wantErr:           true,
		},
		{
			name:              "invalid overflow format",
			domain:            "example.com",
			overflow:          "spf",
			txtMaxSize:        255,
			domainOnRecordKey: false,
			parts:             []string{"v=spf1", "ip4:192.0.2.0/24", "-all"},
			flatten:           []string{},
			want:              nil,
			wantErr:           true,
		},

		{
			name:              "invalid spf record",
			domain:            "example.com",
			overflow:          "spf",
			txtMaxSize:        255,
			domainOnRecordKey: false,
			parts:             []string{"ip4:192.0.2.0/24"},
			flatten:           []string{},
			want:              nil,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := spfbuilder.BuildSPFRecordWithResolver(tt.domain, tt.overflow, tt.txtMaxSize, tt.domainOnRecordKey, tt.parts, tt.flatten, mock)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildSPFRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !mapsEqual(got, tt.want) {
				t.Errorf("BuildSPFRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mapsEqual(a, b map[string][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !slicesEqual(v, bv) {
			return false
		}
	}
	return true
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
