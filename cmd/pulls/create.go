// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdPullsCreate creates a pull request
var CmdPullsCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create a pull-request",
	Description: "Create a pull-request",
	Action:      runPullsCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "head",
			Usage: "Set head branch (default is locally checked out branch)",
		},
		&cli.StringFlag{
			Name:  "head-repo",
			Usage: "Set head repo (default is remote repo for local branch)",
		},
		&cli.StringFlag{
			Name:    "base",
			Aliases: []string{"b"},
			Usage:   "Set base branch (default is default branch)",
		},
	}, flags.IssuePREditFlags...),
}

func runPullsCreate(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{LocalRepo: true})

	// no args -> interactive mode
	if ctx.NumFlags() == 0 {
		return interact.CreatePull(ctx.Login, ctx.Owner, ctx.Repo)
	}

	// else use args to create PR
	opts, err := flags.GetIssuePREditFlags(ctx)
	if err != nil {
		return err
	}

	headRepoSlug := ctx.String("head-repo") // may contain `owner` or `owner/repo`
	headOwner, _ := utils.GetOwnerAndRepo(headRepoSlug, headRepoSlug)

	return task.CreatePull(
		ctx.Login,
		ctx.Owner,
		ctx.Repo,
		ctx.String("base"),
		ctx.String("head"),
		headOwner,
		opts,
	)
}
