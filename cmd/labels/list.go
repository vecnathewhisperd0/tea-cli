// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package labels

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/task"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdLabelsList represents a sub command of labels to list labels
var CmdLabelsList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List labels",
	Description: "List labels",
	Action:      RunLabelsList,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "save",
			Aliases: []string{"s"},
			Usage:   "Save all the labels as a file",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunLabelsList list labels.
func RunLabelsList(cmd *cli.Context) error {
	ctx := config.InitCommand(cmd)
	client := ctx.Login.Client()
	labels, _, err := client.ListRepoLabels(ctx.Owner, ctx.Repo, gitea.ListLabelsOptions{
		ListOptions: flags.GetListOptions(cmd),
	})
	if err != nil {
		log.Fatal(err)
	}

	if ctx.IsSet("save") {
		return task.LabelsExport(labels, ctx.String("save"))
	}

	print.LabelsList(labels, ctx.Output)
	return nil
}
