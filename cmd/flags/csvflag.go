// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package flags

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/modules/utils"
	"github.com/urfave/cli/v2"
)

// CsvFlag is a wrapper around cli.StringFlag, with an added GetValues() method
// to retrieve comma separated string values as a slice.
type CsvFlag struct {
	cli.StringFlag
	AvailableFields []string
}

func NewCsvFlag(name, usage string, aliases, availableFields, defaultFields []string) *CsvFlag {
	return &CsvFlag{
		AvailableFields: availableFields,
		StringFlag: cli.StringFlag{
			Name:    name,
			Aliases: aliases,
			Value:   strings.Join(defaultFields, ","),
			Usage: fmt.Sprintf(`Comma-separated list of %s. Available values:
			%s
		`, usage, strings.Join(availableFields, ",")),
		},
	}
}

func (f CsvFlag) GetValues(ctx *cli.Context) ([]string, error) {
	val := ctx.String(f.Name)
	selection := strings.Split(val, ",")
	if f.AvailableFields != nil && val != "" {
		for _, field := range selection {
			if !utils.Contains(f.AvailableFields, field) {
				return nil, fmt.Errorf("Invalid field '%s'", field)
			}
		}
	}
	return selection, nil
}
