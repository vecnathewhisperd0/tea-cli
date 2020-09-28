// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/modules/intern"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestones represents to operate repositories milestones.
var CmdMilestones = cli.Command{
	Name:        "milestones",
	Aliases:     []string{"ms", "mile"},
	Usage:       "List and create milestones",
	Description: `List and create milestones`,
	ArgsUsage:   "[<milestone name>]",
	Action:      runMilestones,
	Subcommands: []*cli.Command{
		&CmdMilestonesList,
		&CmdMilestonesCreate,
		&CmdMilestonesClose,
		&CmdMilestonesDelete,
		&CmdMilestonesReopen,
		&CmdMilestonesIssues,
	},
	Flags: AllDefaultFlags,
}

func runMilestones(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runMilestoneDetail(ctx.Args().First())
	}
	return runMilestonesList(ctx)
}

func runMilestoneDetail(name string) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)
	client := login.Client()

	milestone, _, err := client.GetMilestoneByName(owner, repo, name)
	if err != nil {
		return err
	}

	print.MilestoneDetails(milestone)
	return nil
}

func editMilestoneStatus(ctx *cli.Context, close bool) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)
	client := login.Client()

	state := gitea.StateOpen
	if close {
		state = gitea.StateClosed
	}
	_, _, err := client.EditMilestoneByName(owner, repo, ctx.Args().First(), gitea.EditMilestoneOption{
		State: &state,
		Title: ctx.Args().First(),
	})

	return err
}
