// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"code.gitea.io/tea/modules/intern"
	"github.com/urfave/cli/v2"
)

// CmdPulls is the main command to operate on PRs
var CmdPulls = cli.Command{
	Name:        "pulls",
	Aliases:     []string{"pull", "pr"},
	Usage:       "List, create, checkout and clean pull requests",
	Description: `List, create, checkout and clean pull requests`,
	ArgsUsage:   "[<pull index>]",
	Action:      runPulls,
	Flags:       IssuePRFlags,
	Subcommands: []*cli.Command{
		&CmdPullsList,
		&CmdPullsCheckout,
		&CmdPullsClean,
		&CmdPullsCreate,
	},
}

func runPulls(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runPullDetail(ctx.Args().First())
	}
	return runPullsList(ctx)
}

func runPullDetail(index string) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	idx, err := argToIndex(index)
	if err != nil {
		return err
	}
	pr, _, err := login.Client().GetPullRequest(owner, repo, idx)
	if err != nil {
		return err
	}

	// TODO: use glamour once #181 is merged
	fmt.Printf("#%d %s\n%s created %s\n\n%s\n", pr.Index,
		pr.Title,
		pr.Poster.UserName,
		pr.Created.Format("2006-01-02 15:04:05"),
		pr.Body,
	)
	return nil
}

func argToIndex(arg string) (int64, error) {
	if strings.HasPrefix(arg, "#") {
		arg = arg[1:]
	}
	return strconv.ParseInt(arg, 10, 64)
}
