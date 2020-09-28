// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/tea/modules/intern"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdRepoCreate represents a sub command of repos to create one
var CmdRepoCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create a repository",
	Description: "Create a repository",
	Action:      runRepoCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{""},
			Required: true,
			Usage:    "name of new repo",
		},
		&cli.StringFlag{
			Name:     "owner",
			Aliases:  []string{"O"},
			Required: false,
			Usage:    "name of repo owner",
		},
		&cli.BoolFlag{
			Name:     "private",
			Required: false,
			Value:    false,
			Usage:    "make repo private",
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"desc"},
			Required: false,
			Usage:    "add description to repo",
		},
		&cli.BoolFlag{
			Name:     "init",
			Required: false,
			Value:    false,
			Usage:    "initialize repo",
		},
		&cli.StringFlag{
			Name:     "labels",
			Required: false,
			Usage:    "name of label set to add",
		},
		&cli.StringFlag{
			Name:     "gitignores",
			Aliases:  []string{"git"},
			Required: false,
			Usage:    "list of gitignore templates (need --init)",
		},
		&cli.StringFlag{
			Name:     "license",
			Required: false,
			Usage:    "add license (need --init)",
		},
		&cli.StringFlag{
			Name:     "readme",
			Required: false,
			Usage:    "use readme template (need --init)",
		},
		&cli.StringFlag{
			Name:     "branch",
			Required: false,
			Usage:    "use custom default branch (need --init)",
		},
	}, LoginOutputFlags...),
}

func runRepoCreate(ctx *cli.Context) error {
	login := intern.InitCommandLoginOnly(globalLoginValue)
	client := login.Client()
	var (
		repo *gitea.Repository
		err  error
	)
	opts := gitea.CreateRepoOption{
		Name:          ctx.String("name"),
		Description:   ctx.String("description"),
		Private:       ctx.Bool("private"),
		AutoInit:      ctx.Bool("init"),
		IssueLabels:   ctx.String("labels"),
		Gitignores:    ctx.String("gitignores"),
		License:       ctx.String("license"),
		Readme:        ctx.String("readme"),
		DefaultBranch: ctx.String("branch"),
	}
	if len(ctx.String("owner")) != 0 {
		repo, _, err = client.CreateOrgRepo(ctx.String("owner"), opts)
	} else {
		repo, _, err = client.CreateRepo(opts)
	}
	if err != nil {
		return err
	}
	if err = runRepoDetail(repo.FullName); err != nil {
		return err
	}
	fmt.Printf("%s\n", repo.HTMLURL)
	return nil
}
