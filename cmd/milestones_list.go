// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/tea/modules/intern"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestonesList represents a sub command of milestones to list milestones
var CmdMilestonesList = cli.Command{
	Name:        "ls",
	Usage:       "List milestones of the repository",
	Description: `List milestones of the repository`,
	Action:      runMilestonesList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by milestone state (all|open|closed)",
			DefaultText: "open",
		},
		&PaginationPageFlag,
		&PaginationLimitFlag,
	}, AllDefaultFlags...),
}

func runMilestonesList(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "closed":
		state = gitea.StateClosed
	}

	milestones, _, err := login.Client().ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{
		ListOptions: getListOptions(ctx),
		State:       state,
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Title",
	}
	if state == gitea.StateAll {
		headers = append(headers, "State")
	}
	headers = append(headers,
		"Open/Closed Issues",
		"DueDate",
	)

	var values [][]string

	for _, m := range milestones {
		var deadline = ""

		if m.Deadline != nil && !m.Deadline.IsZero() {
			deadline = m.Deadline.Format("2006-01-02 15:04:05")
		}

		item := []string{
			m.Title,
		}
		if state == gitea.StateAll {
			item = append(item, string(m.State))
		}
		item = append(item,
			fmt.Sprintf("%d/%d", m.OpenIssues, m.ClosedIssues),
			deadline,
		)

		values = append(values, item)
	}
	Output(globalOutputValue, headers, values)

	return nil
}
