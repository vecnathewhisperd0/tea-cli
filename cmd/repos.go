// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli/v2"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "Operate with repositories",
	Description: `Operate with repositories`,
	Action:      runReposList,
	Subcommands: []*cli.Command{
		&CmdReposList,
	},
	Flags: LoginOutputFlags,
}

// CmdReposList represents a sub command of issues to list issues
var CmdReposList = cli.Command{
	Name:        "ls",
	Usage:       "List available repositories",
	Description: `List available repositories`,
	Action:      runReposList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "mode",
			Aliases:  []string{"m"},
			Required: false,
			Usage:    "Filter by mode: fork, mirror, source",
		},
		&cli.StringFlag{
			Name:     "org",
			Required: false,
			Usage:    "Filter by organization",
		},
		&cli.StringFlag{
			Name:     "user",
			Aliases:  []string{"u"},
			Required: false,
			Usage:    "Filter by user",
		},
	}, LoginOutputFlags...),
}

// runReposList list repositories
func runReposList(ctx *cli.Context) error {
	login := initCommandLoginOnly()

	mode := ctx.String("mode")
	org := ctx.String("org")
	user := ctx.String("user")

	var rps []*gitea.Repository
	var err error

	// ToDo: on sdk v0.13.0 release, switch to SearchRepos()
	// Note: user filter can be used as org filter too
	if org != "" {
		rps, err = login.Client().ListOrgRepos(org, gitea.ListOrgReposOptions{})
	} else if user != "" {
		rps, err = login.Client().ListUserRepos(user, gitea.ListReposOptions{})
	} else {
		rps, err = login.Client().ListMyRepos(gitea.ListReposOptions{})
	}
	if err != nil {
		log.Fatal(err)
	}

	var repos []*gitea.Repository
	if mode == "" {
		repos = rps
	} else if mode == "fork" {
		for _, rp := range rps {
			if rp.Fork == true {
				repos = append(repos, rp)
			}
		}
	} else if mode == "mirror" {
		for _, rp := range rps {
			if rp.Mirror == true {
				repos = append(repos, rp)
			}
		}
	} else if mode == "source" {
		for _, rp := range rps {
			if rp.Mirror != true && rp.Fork != true {
				repos = append(repos, rp)
			}
		}
	} else {
		log.Fatal("Unknown mode: ", mode, "\nUse one of the following:\n- fork\n- mirror\n- source\n")
		return nil
	}

	if len(rps) == 0 {
		log.Fatal("No repositories found", rps)
		return nil
	}

	headers := []string{
		"Name",
		"Type",
		"SSH",
		"Owner",
	}
	var values [][]string

	for _, rp := range repos {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}

		values = append(
			values,
			[]string{
				rp.FullName,
				mode,
				rp.SSHURL,
				rp.Owner.UserName,
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}
