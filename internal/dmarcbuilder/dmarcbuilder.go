// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package dmarcbuilder

import (
	"errors"
	"fmt"
	"strings"
)

type DMARCConfig struct {
	Version         string
	Policy          string
	SubdomainPolicy string
	AlignmentSPF    string
	AlignmentDKIM   string
	Percent         int32
	RUA             []string
	RUF             []string
	FailureOptions  string
	FailureFormat   string
	ReportInterval  int32
}

func DmarcBuilder(value DMARCConfig) (string, error) {
	if value.Version == "" {
		value.Version = "DMARC1"
	}

	if value.Policy == "" {
		value.Policy = "none"
	}

	validPolicies := map[string]bool{"none": true, "quarantine": true, "reject": true}
	if !validPolicies[value.Policy] {
		return "", errors.New("invalid DMARC policy")
	}

	record := []string{"v=" + value.Version, "p=" + value.Policy}

	if value.SubdomainPolicy != "" {
		if !validPolicies[value.SubdomainPolicy] {
			return "", errors.New("invalid DMARC subdomain policy")
		}
		record = append(record, "sp="+value.SubdomainPolicy)
	}

	alignments := map[string]string{"relaxed": "r", "strict": "s", "r": "r", "s": "s"}
	if val, ok := alignments[value.AlignmentDKIM]; ok {
		record = append(record, "adkim="+val)
	} else if value.AlignmentDKIM != "" {
		return "", errors.New("invalid DMARC DKIM alignment policy")
	}

	if val, ok := alignments[value.AlignmentSPF]; ok {
		record = append(record, "aspf="+val)
	} else if value.AlignmentSPF != "" {
		return "", errors.New("invalid DMARC SPF alignment policy")
	}

	if value.Percent > 0 {
		record = append(record, fmt.Sprintf("pct=%d", value.Percent))
	}

	if len(value.RUA) > 0 {
		record = append(record, "rua="+strings.Join(value.RUA, ","))
	}

	if len(value.RUF) > 0 {
		record = append(record, "ruf="+strings.Join(value.RUF, ","))
	}

	if len(value.RUF) > 0 && value.FailureOptions != "" {
		validFailureOptions := map[string]bool{"0": true, "1": true, "d": true, "s": true}
		if validFailureOptions[value.FailureOptions] {
			record = append(record, "fo="+value.FailureOptions)
		} else {
			return "", errors.New("invalid DMARC failure options")
		}
	}

	if len(value.RUF) > 0 && value.FailureFormat != "" {
		record = append(record, "rf="+value.FailureFormat)
	}

	if value.ReportInterval > 0 {
		record = append(record, fmt.Sprintf("ri=%d", value.ReportInterval))
	}

	return strings.Join(record, "; "), nil
}
