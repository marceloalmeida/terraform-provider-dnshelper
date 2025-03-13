// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package caabuilder_test

import (
	"testing"

	"github.com/marceloalmeida/terraform-provider-dnshelper/internal/caabuilder"
)

func TestCAABuilder_NoIssue(t *testing.T) {
	tests := []struct {
		name    string
		args    caabuilder.CAAConfig
		want    []caabuilder.CAARecord
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := caabuilder.CAABuilder(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CAABuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("CAABuilder() = %v, want %v", got, tt.want)
				}
				for i := range got {
					if got[i].Tag != tt.want[i].Tag {
						t.Errorf("CAABuilder() = %v, want %v", got, tt.want)
					}
					if got[i].Value != tt.want[i].Value {
						t.Errorf("CAABuilder() = %v, want %v", got, tt.want)
					}
					if got[i].Flag != tt.want[i].Flag {
						t.Errorf("CAABuilder() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
func TestCAABuilderString(t *testing.T) {
	tests := []struct {
		name    string
		args    caabuilder.CAAConfig
		want    []string
		wantErr bool
	}{
		{
			name:    "Empty Config",
			args:    caabuilder.CAAConfig{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Single Issue",
			args: caabuilder.CAAConfig{
				Issue: []string{"letsencrypt.org"},
			},
			want: []string{
				`0 issue "letsencrypt.org"`,
			},
			wantErr: false,
		},
		{
			name: "Multiple Records",
			args: caabuilder.CAAConfig{
				Iodef:             "mailto:security@example.com",
				IodefCritical:     true,
				Issue:             []string{"letsencrypt.org"},
				Issuewild:         []string{"sectigo.com"},
				IssuewildCritical: true,
			},
			want: []string{
				`128 iodef "mailto:security@example.com"`,
				`0 issue "letsencrypt.org"`,
				`128 issuewild "sectigo.com"`,
			},
			wantErr: false,
		},

		{
			name: "Multiple Records with IodefCritical false",
			args: caabuilder.CAAConfig{
				Iodef:             "mailto:security@example.com",
				IodefCritical:     false,
				Issue:             []string{"letsencrypt.org"},
				Issuewild:         []string{"sectigo.com"},
				IssueCritical:     false,
				IssuewildCritical: true,
			},
			want: []string{
				`0 iodef "mailto:security@example.com"`,
				`0 issue "letsencrypt.org"`,
				`128 issuewild "sectigo.com"`,
			},
			wantErr: false,
		},
		{
			name: "Multiple Records with IodefCritical false and IssueCritical true",
			args: caabuilder.CAAConfig{
				Iodef:             "mailto:security@example.com",
				IodefCritical:     false,
				Issue:             []string{"letsencrypt.org"},
				Issuewild:         []string{"sectigo.com"},
				IssueCritical:     true,
				IssuewildCritical: true,
			},
			want: []string{
				`0 iodef "mailto:security@example.com"`,
				`128 issue "letsencrypt.org"`,
				`128 issuewild "sectigo.com"`,
			},
			wantErr: false,
		},
		{
			name: "None Values",
			args: caabuilder.CAAConfig{
				Issue:     []string{"none"},
				Issuewild: []string{"none"},
			},
			want: []string{
				`0 issue ";"`,
				`0 issuewild ";"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := caabuilder.CAABuilderString(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CAABuilderString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("CAABuilderString() = %v, want %v", got, tt.want)
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("CAABuilderString() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
