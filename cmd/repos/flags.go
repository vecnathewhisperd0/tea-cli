// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

// printFieldsFlag provides a selection of fields to print
var printFieldsFlag = cli.StringFlag{
	Name:    "fields",
	Aliases: []string{"f"},
	Usage: fmt.Sprintf(`Comma-separated list of fields to print. Available values:
		%v
	 `, print.RepoFields),
	Value: "url,stars,forks,description",
}

func getFields(ctx *cli.Context) []string {
	return strings.Split(ctx.String("fields"), ",")
}
