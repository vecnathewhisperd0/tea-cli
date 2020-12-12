// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package labels

import (
	"log"

	"code.gitea.io/tea/modules/config"

	"github.com/urfave/cli/v2"
)

// CmdLabelDelete represents a sub command of labels to delete label.
var CmdLabelDelete = cli.Command{
	Name:        "delete",
	Usage:       "Delete a label",
	Description: `Delete a label`,
	Action:      runLabelDelete,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "label id",
		},
	},
}

func runLabelDelete(cmd *cli.Context) error {
	ctx := config.InitCommand(cmd)
	ctx.Ensure(config.CtxRequirement{RemoteRepo: true})

	_, err := ctx.Login.Client().DeleteLabel(ctx.Owner, ctx.Repo, ctx.Int64("id"))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
