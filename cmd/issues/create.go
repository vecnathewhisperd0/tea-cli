// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

// CmdIssuesCreate represents a sub command of issues to create issue
var CmdIssuesCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create an issue on repository",
	Description: `Create an issue on repository`,
	Action:      runIssuesCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "issue title to create",
		},
		&cli.StringFlag{
			Name:    "body",
			Aliases: []string{"b"},
			Usage:   "issue body to create",
		},
		&cli.StringFlag{
			Name:    "assignees",
			Aliases: []string{"a"},
			Usage:   "Comma separated list of usernames to assign",
		},
		&cli.StringFlag{
			Name:    "deadline",
			Aliases: []string{"D"},
			Usage:   "Deadline timestamp to assign",
		},
		&cli.StringFlag{
			Name:    "labels",
			Aliases: []string{"L"},
			Usage:   "Comma separated list of labels to assign",
		},
		&cli.StringFlag{
			Name:    "milestone",
			Aliases: []string{"m"},
			Usage:   "Milestone to assign",
		},
	}, flags.LoginRepoFlags...),
}

func runIssuesCreate(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if ctx.NumFlags() == 0 {
		return interact.CreateIssue(ctx.Login, ctx.Owner, ctx.Repo)
	}

	var (
		client      *gitea.Client
		milestoneID int64
		deadline    *time.Time
	)

	date := ctx.String("deadline")
	if date != "" {
		t, err := dateparse.ParseAny(date)
		if err != nil {
			return err
		}
		deadline = &t
	}

	labelNames := strings.Split(ctx.String("labels"), ",")
	labelIDs := make([]int64, len(labelNames))
	if len(labelNames) != 0 {
		if client == nil {
			client = ctx.Login.Client()
		}
		labels, _, err := client.ListRepoLabels(ctx.Owner, ctx.Repo, gitea.ListLabelsOptions{})
		if err != nil {
			return err
		}
		for _, l := range labels {
			if utils.Contains(labelNames, l.Name) {
				labelIDs = append(labelIDs, l.ID)
			}
		}
	}

	if milestoneName := ctx.String("milestone"); len(milestoneName) != 0 {
		if client == nil {
			client = ctx.Login.Client()
		}
		ms, _, err := client.GetMilestoneByName(ctx.Owner, ctx.Repo, milestoneName)
		if err != nil {
			return fmt.Errorf("Milestone '%s' not found", milestoneName)
		}
		milestoneID = ms.ID
	}

	return task.CreateIssue(
		ctx.Login,
		ctx.Owner,
		ctx.Repo,
		gitea.CreateIssueOption{
			Title:     ctx.String("title"),
			Body:      ctx.String("body"),
			Assignees: strings.Split(ctx.String("assignees"), ","),
			Deadline:  deadline,
			Labels:    labelIDs,
			Milestone: milestoneID,
		},
	)
}
