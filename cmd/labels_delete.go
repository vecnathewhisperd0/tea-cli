// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/tea/modules/intern"

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

func runLabelDelete(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	_, err := login.Client().DeleteLabel(owner, repo, ctx.Int64("id"))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
