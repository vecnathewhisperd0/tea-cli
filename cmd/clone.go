// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"

	"github.com/urfave/cli/v2"
)

// CmdRepoClone represents a sub command of repos to create a local copy
var CmdRepoClone = cli.Command{
	Name:        "clone",
	Aliases:     []string{"C"},
	Usage:       "Clone a repository locally",
	Description: "Clone a repository locally, without a local git installation required",
	Category:    catHelpers,
	Action:      runRepoClone,
	ArgsUsage:   "[target dir]",
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:    "depth",
			Aliases: []string{"d"},
			Usage:   "num commits to fetch, defaults to all",
		},
	}, flags.LoginRepoFlags...),
}

func runRepoClone(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	_, err := task.RepoClone(
		ctx.Args().First(),
		ctx.Login,
		ctx.Owner,
		ctx.Repo,
		interact.PromptPassword,
		ctx.Int("depth"),
	)

	return err
}
