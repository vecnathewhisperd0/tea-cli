// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package flags

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/task"

	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

// StateFlag provides flag to specify issue/pr state, defaulting to "open"
var StateFlag = cli.StringFlag{
	Name:        "state",
	Usage:       "Filter by state (all|open|closed)",
	DefaultText: "open",
}

// IssuePRFlags defines flags that should be available on issue & pr listing flags.
var IssuePRFlags = append([]cli.Flag{
	&StateFlag,
	&PaginationPageFlag,
	&PaginationLimitFlag,
}, AllDefaultFlags...)

// IssuePREditFlags defines flags for properties of issues and PRs
var IssuePREditFlags = append([]cli.Flag{
	&cli.StringFlag{
		Name:    "title",
		Aliases: []string{"t"},
	},
	&cli.StringFlag{
		Name:    "description",
		Aliases: []string{"d"},
	},
	&cli.StringFlag{
		Name:    "assignees",
		Aliases: []string{"a"},
		Usage:   "Comma-separated list of usernames to assign",
	},
	&cli.StringFlag{
		Name:    "labels",
		Aliases: []string{"L"},
		Usage:   "Comma-separated list of labels to assign",
	},
	&cli.StringFlag{
		Name:    "deadline",
		Aliases: []string{"D"},
		Usage:   "Deadline timestamp to assign",
	},
	&cli.StringFlag{
		Name:    "milestone",
		Aliases: []string{"m"},
		Usage:   "Milestone to assign",
	},
}, LoginRepoFlags...)

// GetIssuePREditFlags parses all IssuePREditFlags
func GetIssuePREditFlags(ctx *context.TeaContext) (*gitea.CreateIssueOption, error) {
	opts := gitea.CreateIssueOption{
		Title:     ctx.String("title"),
		Body:      ctx.String("description"),
		Assignees: strings.Split(ctx.String("assignees"), ","),
	}
	var err error

	date := ctx.String("deadline")
	if date != "" {
		t, err := dateparse.ParseAny(date)
		if err != nil {
			return nil, err
		}
		opts.Deadline = &t
	}

	client := ctx.Login.Client()

	labelNames := strings.Split(ctx.String("labels"), ",")
	if len(labelNames) != 0 {
		if client == nil {
			client = ctx.Login.Client()
		}
		if opts.Labels, err = task.ResolveLabelNames(client, ctx.Owner, ctx.Repo, labelNames); err != nil {
			return nil, err
		}
	}

	if milestoneName := ctx.String("milestone"); len(milestoneName) != 0 {
		if client == nil {
			client = ctx.Login.Client()
		}
		ms, _, err := client.GetMilestoneByName(ctx.Owner, ctx.Repo, milestoneName)
		if err != nil {
			return nil, fmt.Errorf("Milestone '%s' not found", milestoneName)
		}
		opts.Milestone = ms.ID
	}

	return &opts, nil
}
