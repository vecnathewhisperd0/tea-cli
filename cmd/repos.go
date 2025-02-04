// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"code.gitea.io/tea/cmd/repos"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Aliases:     []string{"repo"},
	Category:    catEntities,
	Usage:       "Show repository details",
	Description: "Show repository details",
	ArgsUsage:   "[<repo owner>/<repo name>]",
	Action:      runRepos,
	Subcommands: []*cli.Command{
		&repos.CmdReposList,
		&repos.CmdReposSearch,
		&repos.CmdRepoCreate,
		&repos.CmdRepoCreateFromTemplate,
		&repos.CmdRepoFork,
		&repos.CmdRepoMigrate,
		&repos.CmdRepoRm,
	},
	Flags: repos.CmdReposListFlags,
}

func runRepos(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runRepoDetail(ctx, ctx.Args().First())
	}
	return repos.RunReposList(ctx)
}

func runRepoDetail(cmd *cli.Context, path string) error {
	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()
	repoOwner, repoName := utils.GetOwnerAndRepo(path, ctx.Owner)
	repo, _, err := client.GetRepo(repoOwner, repoName)
	if err != nil {
		return err
	}
	topics, _, err := client.ListRepoTopics(repoOwner, repoName, gitea.ListRepoTopicsOptions{})
	if err != nil {
		return err
	}

	print.RepoDetails(repo, topics)
	return nil
}
