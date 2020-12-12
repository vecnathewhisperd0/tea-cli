// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdPullsCheckout is a command to locally checkout the given PR
var CmdPullsCheckout = cli.Command{
	Name:        "checkout",
	Usage:       "Locally check out the given PR",
	Description: `Locally check out the given PR`,
	Action:      runPullsCheckout,
	ArgsUsage:   "<pull index>",
	Flags:       flags.AllDefaultFlags,
}

func runPullsCheckout(cmd *cli.Context) error {
	ctx := config.InitCommand(cmd)
	if ctx.Args().Len() != 1 {
		log.Fatal("Must specify a PR index")
	}
	idx, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	return task.PullCheckout(ctx.Login, ctx.Owner, ctx.Repo, idx, interact.PromptPassword)
}
