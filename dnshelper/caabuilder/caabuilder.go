// Copyright (c) Marcelo Almeida
// SPDX-License-Identifier: MPL-2.0

package caabuilder

import (
	"errors"
	"strconv"
)

type CAARecord struct {
	Tag   string
	Value string
	Flag  int
}

const caaCritical = 128

func CAA(tag string, value string, flag int) CAARecord {
	return CAARecord{
		Tag:   tag,
		Value: value,
		Flag:  flag,
	}
}

type CAAConfig struct {
	Iodef             string
	IodefCritical     bool
	Issue             []string
	IssueCritical     bool
	Issuewild         []string
	IssuewildCritical bool
}

func CAABuilder(value CAAConfig) ([]CAARecord, error) {
	if len(value.Issue) == 1 && value.Issue[0] == "none" {
		value.Issue = []string{";"}
	}
	if len(value.Issuewild) == 1 && value.Issuewild[0] == "none" {
		value.Issuewild = []string{";"}
	}

	if len(value.Issue) == 0 && len(value.Issuewild) == 0 {
		return nil, errors.New("CAABuilder requires at least one entry in issue or issuewild")
	}

	r := []CAARecord{}

	if value.Iodef != "" {
		if value.IodefCritical {
			r = append(r, CAA("iodef", value.Iodef, caaCritical))
		} else {
			r = append(r, CAA("iodef", value.Iodef, 0))
		}
	}

	if len(value.Issue) > 0 {
		flag := 0
		if value.IssueCritical {
			flag = caaCritical
		}
		for _, issue := range value.Issue {
			r = append(r, CAA("issue", issue, flag))
		}
	}

	if len(value.Issuewild) > 0 {
		flag := 0
		if value.IssuewildCritical {
			flag = caaCritical
		}
		for _, issuewild := range value.Issuewild {
			r = append(r, CAA("issuewild", issuewild, flag))
		}
	}

	return r, nil
}

func CAABuilderString(value CAAConfig) ([]string, error) {
	records, err := CAABuilder(value)
	if err != nil {
		return nil, err
	}

	r := []string{}
	for _, record := range records {
		r = append(r, strconv.Itoa(record.Flag)+" "+record.Tag+" "+`"`+record.Value+`"`)
	}

	return r, nil
}
