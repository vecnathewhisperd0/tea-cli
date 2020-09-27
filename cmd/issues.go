// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/tea/modules/intern"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdIssues represents to login a gitea server.
var CmdIssues = cli.Command{
	Name:        "issues",
	Usage:       "List, create and update issues",
	Description: "List, create and update issues",
	ArgsUsage:   "[<issue index>]",
	Action:      runIssues,
	Subcommands: []*cli.Command{
		&CmdIssuesList,
		&CmdIssuesCreate,
		&CmdIssuesReopen,
		&CmdIssuesClose,
	},
	Flags: IssuePRFlags,
}

func runIssues(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runIssueDetail(ctx.Args().First())
	}
	return runIssuesList(ctx)
}

func runIssueDetail(index string) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	idx, err := argToIndex(index)
	if err != nil {
		return err
	}
	issue, _, err := login.Client().GetIssue(owner, repo, idx)
	if err != nil {
		return err
	}
	print.IssueDetails(issue)
	return nil
}

// editIssueState abstracts the arg parsing to edit the given issue
func editIssueState(ctx *cli.Context, opts gitea.EditIssueOption) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)
	if ctx.Args().Len() == 0 {
		log.Fatal(ctx.Command.ArgsUsage)
	}

	index, err := argToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	_, _, err = login.Client().EditIssue(owner, repo, index, opts)
	// TODO: print (short)IssueDetails
	return err
}
