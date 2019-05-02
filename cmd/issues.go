// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdIssues represents to login a gitea server.
var CmdIssues = cli.Command{
	Name:        "issues",
	Usage:       "Operate with issues of the repository",
	Description: `Operate with issues of the repository`,
	Action:      runIssues,
	Subcommands: []cli.Command{
		CmdIssuesList,
		CmdIssuesCreate,
	},
	Flags: append([]cli.Flag{}, RepoDefaultFlags...),
}

// CmdIssuesList represents a sub command of issues to list issues
var CmdIssuesList = cli.Command{
	Name:        "ls",
	Usage:       "List issues of the repository",
	Description: `List issues of the repository`,
	Action:      runIssuesList,
}

func runIssues(ctx *cli.Context) error {
	if len(os.Args) == 3 {
		return runIssueDetail(ctx, os.Args[2])
	}
	return runIssuesList(ctx)
}

func runIssueDetail(ctx *cli.Context, index string) error {
	login, owner, repo := initCommand()

	if strings.HasPrefix(index, "#") {
		index = index[1:]
	}

	idx, err := strconv.ParseInt(index, 10, 64)
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

	issues, err := login.Client().ListRepoIssues(owner, repo, gitea.ListIssueOption{
		Page:  0,
		State: string(gitea.StateOpen),
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(issues) == 0 {
		fmt.Println("No issues left")
		return nil
	}

	for _, issue := range issues {
		name := issue.Poster.FullName
		if len(name) == 0 {
			name = issue.Poster.UserName
		}
		fmt.Printf("#%d\t%s\t%s\t%s\n", issue.Index, name, issue.Updated.Format("2006-01-02 15:04:05"), issue.Title)
	}

	return nil
}

// CmdIssuesCreate represents a sub command of issues to create issue
var CmdIssuesCreate = cli.Command{
	Name:        "create",
	Usage:       "Create an issue on repository",
	Description: `Create an issue on repository`,
	Action:      runIssuesCreate,
	Flags: append([]cli.Flag{
		cli.StringFlag{
			Name:  "title, t",
			Usage: "issue title to create",
		},
		cli.StringFlag{
			Name:  "body, b",
			Usage: "issue body to create",
		},
	}, RepoDefaultFlags...),
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
