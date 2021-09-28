// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"net/url"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/interact"

	"github.com/urfave/cli/v2"
)

// CmdRepoClone represents a sub command of repos to create a local copy
var CmdRepoClone = cli.Command{
	Name:        "clone",
	Aliases:     []string{"C"},
	Usage:       "Clone a repository locally",
	Description: "Clone a repository locally, without a local git installation required (defaults to PWD)",
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

	// get clone URLs robustly
	repo, _, err := ctx.Login.Client().GetRepo(ctx.Owner, ctx.Repo)
	if err != nil {
		return err
	}
	var url *url.URL
	urlStr := repo.CloneURL
	if ctx.Login.SSHKey != "" {
		urlStr = repo.SSHURL
	}
	url, err = git.ParseURL(urlStr)
	if err != nil {
		return err
	}

	auth, err := git.GetAuthForURL(url, ctx.Login.Token, ctx.Login.SSHKey, interact.PromptPassword)
	if err != nil {
		return err
	}

	// default path behaviour as native git
	localPath := ctx.Args().First()
	if localPath == "" {
		localPath = ctx.Repo
	}

	_, err = git.CloneIntoPath(
		url.String(),
		localPath,
		auth,
		ctx.Int("depth"),
		ctx.Login.Insecure,
	)

	return err
}
