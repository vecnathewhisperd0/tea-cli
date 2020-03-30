// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"strconv"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli/v2"
)

// CmdIssues represents to login a gitea server.
var CmdIssues = cli.Command{
	Name:        "issues",
	Usage:       "List and create issues",
	Description: `List and create issues`,
	Action:      runIssues,
	Subcommands: []*cli.Command{
		&CmdIssuesList,
		&CmdIssuesCreate,
	},
	Flags: AllDefaultFlags,
}

// CmdIssuesList represents a sub command of issues to list issues
var CmdIssuesList = cli.Command{
	Name:        "ls",
	Usage:       "List issues of the repository",
	Description: `List issues of the repository`,
	Action:      runIssuesList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by issue state (all|open|closed)",
			DefaultText: "open",
		},
	}, AllDefaultFlags...),
}

func runIssues(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runIssueDetail(ctx, ctx.Args().First())
	}
	return runIssuesList(ctx)
}

func runIssueDetail(ctx *cli.Context, index string) error {
	login, owner, repo := initCommand()

	idx, err := argToIndex(index)
	if err != nil {
		return err
	}
	issue, err := login.Client().GetIssue(owner, repo, idx)
	if err != nil {
		return err
	}

	fmt.Printf("#%d %s\n%s created %s\n\n%s", issue.Index,
		issue.Title,
		issue.Poster.UserName,
		issue.Created.Format("2006-01-02 15:04:05"),
		issue.Body,
	)
	return nil
}

func runIssuesList(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "open":
		state = gitea.StateOpen
	case "closed":
		state = gitea.StateClosed
	}

	issues, err := login.Client().ListRepoIssues(owner, repo, gitea.ListIssueOption{
		Page:  0,
		State: string(state),
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Index",
		"State",
		"Author",
		"Updated",
		"Title",
	}

	var values [][]string

	if len(issues) == 0 {
		Output(outputValue, headers, values)
		return nil
	}

	for _, issue := range issues {
		name := issue.Poster.FullName
		if len(name) == 0 {
			name = issue.Poster.UserName
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(issue.Index, 10),
				string(issue.State),
				name,
				issue.Updated.Format("2006-01-02 15:04:05"),
				issue.Title,
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}

// CmdIssuesCreate represents a sub command of issues to create issue
var CmdIssuesCreate = cli.Command{
	Name:        "create",
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
	}, LoginRepoFlags...),
}

func runIssuesCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	_, err := login.Client().CreateIssue(owner, repo, gitea.CreateIssueOption{
		Title: ctx.String("title"),
		Body:  ctx.String("body"),
		// TODO:
		//Assignee  string   `json:"assignee"`
		//Assignees []string `json:"assignees"`
		//Deadline *time.Time `json:"due_date"`
		//Milestone int64 `json:"milestone"`
		//Labels []int64 `json:"labels"`
		//Closed bool    `json:"closed"`
	})

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
