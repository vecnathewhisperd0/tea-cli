// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "Operate with repositories",
	Description: `Operate with repositories`,
	Action:      runReposList,
	Subcommands: []cli.Command{
		CmdReposList,
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
		cli.StringFlag{
			Name:  "mode",
			Usage: "Filter listed repositories based on mode, optional - fork, mirror, source",
		},
		cli.StringFlag{
			Name:  "org",
			Usage: "Filter listed repositories based on organization, optional",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Filter listed repositories absed on user, optional",
		},
	}, LoginRepoFlags...),
}

// runReposList list repositories
func runReposList(ctx *cli.Context) error {
	login := initCommandLoginOnly()

	mode := ctx.String("mode")
	org := ctx.String("org")
	user := ctx.String("user")

	var rps []*gitea.Repository
	var err error

	if org != "" {
		rps, err = login.Client().ListOrgRepos(org)
	} else if user != "" {
		rps, err = login.Client().ListUserRepos(user)
	} else {
		rps, err = login.Client().ListMyRepos()
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
		fmt.Printf("Unknown mode '%s'\nUse one of the following:\n- fork\n- mirror\n- source\n", mode)
		return nil
	}

	if len(rps) == 0 {
		fmt.Println("No repositories found")
		return nil
	}

	fmt.Println("Name | Type/Mode | SSH-URL | Owner")
	for _, rp := range repos {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}
		fmt.Printf("%s | %s | %s | %s\n", rp.FullName, mode, rp.SSHURL, rp.Owner.UserName)
	}

	return nil
}
