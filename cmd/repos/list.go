// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReposListFlags contains all flags needed for repo listing
var CmdReposListFlags = append([]cli.Flag{
	&cli.BoolFlag{
		Name:     "watched",
		Aliases:  []string{"w"},
		Required: false,
		Usage:    "List your watched repos instead",
	},
	&flags.PaginationPageFlag,
	&flags.PaginationLimitFlag,
}, flags.LoginOutputFlags...)

// CmdReposList represents a sub command of repos to list them
var CmdReposList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List repositories you have access to",
	Description: "List repositories you have access to",
	Action:      RunReposList,
	Flags:       CmdReposListFlags,
}

// RunReposList list repositories
func RunReposList(ctx *cli.Context) error {
	login := config.InitCommandLoginOnly(flags.GlobalLoginValue)
	client := login.Client()

	var rps []*gitea.Repository
	var err error
	if ctx.Bool("watched") {
		rps, _, err = client.GetMyWatchedRepos() // TODO: this does not expose pagination..
	} else {
		rps, _, err = client.ListMyRepos(gitea.ListReposOptions{
			ListOptions: flags.GetListOptions(ctx),
		})
	}

	if err != nil {
		return err
	}

	print.ReposList(rps)
	return nil
}
