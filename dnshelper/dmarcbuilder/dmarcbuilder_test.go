// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package dmarcbuilder_test

import (
	"testing"

	"github.com/marceloalmeida/terraform-provider-dnshelper/dnshelper/dmarcbuilder"
)

func TestDmarcBuilder(t *testing.T) {
	tests := []struct {
		name    string
		args    dmarcbuilder.DMARCConfig
		want    string
		wantErr bool
	}{
		{
			name: "Invalid Failure Options",
			args: dmarcbuilder.DMARCConfig{
				RUF:            []string{"mailto:forensics@example.com"},
				FailureOptions: "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Failure Format",
			args: dmarcbuilder.DMARCConfig{
				RUF:           []string{"mailto:forensics@example.com"},
				FailureFormat: "afrf",
			},
			want:    "v=DMARC1; p=none; ruf=mailto:forensics@example.com; rf=afrf",
			wantErr: false,
		},
		{
			name:    "Basic DMARC Record",
			args:    dmarcbuilder.DMARCConfig{},
			want:    "v=DMARC1; p=none",
			wantErr: false,
		},
		{
			name: "Invalid Policy",
			args: dmarcbuilder.DMARCConfig{
				Policy: "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid Subdomain Policy",
			args: dmarcbuilder.DMARCConfig{
				SubdomainPolicy: "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Valid Alignments",
			args: dmarcbuilder.DMARCConfig{
				AlignmentDKIM: "strict",
				AlignmentSPF:  "relaxed",
			},
			want:    "v=DMARC1; p=none; adkim=s; aspf=r",
			wantErr: false,
		},
		{
			name: "Invalid DKIM Alignment",
			args: dmarcbuilder.DMARCConfig{
				AlignmentDKIM: "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid SPF Alignment",
			args: dmarcbuilder.DMARCConfig{
				AlignmentSPF: "invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Percentage Setting",
			args: dmarcbuilder.DMARCConfig{
				Percent: 50,
			},
			want:    "v=DMARC1; p=none; pct=50",
			wantErr: false,
		},
		{
			name: "RUA Settings",
			args: dmarcbuilder.DMARCConfig{
				RUA: []string{"mailto:dmarc@example.com"},
			},
			want:    "v=DMARC1; p=none; rua=mailto:dmarc@example.com",
			wantErr: false,
		},
		{
			name: "RUF Settings",
			args: dmarcbuilder.DMARCConfig{
				RUF: []string{"mailto:forensics@example.com"},
			},
			want:    "v=DMARC1; p=none; ruf=mailto:forensics@example.com",
			wantErr: false,
		},
		{
			name: "Complete DMARC Record",
			args: dmarcbuilder.DMARCConfig{
				Version:         "DMARC1",
				Policy:          "reject",
				SubdomainPolicy: "quarantine",
				AlignmentDKIM:   "strict",
				AlignmentSPF:    "relaxed",
				Percent:         100,
				RUA:             []string{"mailto:dmarc@example.com"},
				RUF:             []string{"mailto:forensics@example.com"},
				ReportInterval:  86400,
			},
			want:    "v=DMARC1; p=reject; sp=quarantine; adkim=s; aspf=r; pct=100; rua=mailto:dmarc@example.com; ruf=mailto:forensics@example.com; ri=86400",
			wantErr: false,
		},
		{
			name: "Multiple Reporting Addresses",
			args: dmarcbuilder.DMARCConfig{
				RUA: []string{"mailto:dmarc1@example.com", "mailto:dmarc2@example.com"},
				RUF: []string{"mailto:forensics1@example.com", "mailto:forensics2@example.com"},
			},
			want:    "v=DMARC1; p=none; rua=mailto:dmarc1@example.com,mailto:dmarc2@example.com; ruf=mailto:forensics1@example.com,mailto:forensics2@example.com",
			wantErr: false,
		},
		{
			name: "Failure Options 1",
			args: dmarcbuilder.DMARCConfig{
				RUF:            []string{"mailto:forensics@example.com"},
				FailureOptions: "1",
			},
			want:    "v=DMARC1; p=none; ruf=mailto:forensics@example.com; fo=1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dmarcbuilder.DmarcBuilder(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DmarcBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DmarcBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}
