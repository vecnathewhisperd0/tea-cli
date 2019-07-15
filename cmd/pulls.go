// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdPulls represents to login a gitea server.
var CmdPulls = cli.Command{
	Name:        "pulls",
	Usage:       "Operate with pulls of the repository",
	Description: `Operate with pulls of the repository`,
	Action:      runPulls,
	Subcommands: []cli.Command{
		CmdPullsCreate,
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login, optional when inside a gitea repository",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "Indicate one repository, optional when inside a gitea repository",
		},
	},
}

func runPulls(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)

	prs, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		Page:  0,
		State: string(gitea.StateOpen),
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(prs) == 0 {
		fmt.Println("No pull requests left")
		return nil
	}

	for _, pr := range prs {
		if pr == nil {
			continue
		}
		name := pr.Poster.FullName
		if len(name) == 0 {
			name = pr.Poster.UserName
		}
		fmt.Printf("#%d\t%s\t%s\t%s\n", pr.Index, name, pr.Updated.Format("2006-01-02 15:04:05"), pr.Title)
	}

	return nil
}

// CmdPullsCreate represents a sub command of pulls to create pr
var CmdPullsCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a pull-request on repository",
	Description: `Create a pull-request on repository`,
	Action:      runCreatePullRequest,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "head",
			Usage: "pull-request head",
		},
		cli.StringFlag{
			Name:  "base, b",
			Usage: "pull-request base",
		},
		cli.StringFlag{
			Name:  "title, t",
			Usage: "pull-request title",
		},
		cli.StringFlag{
			Name:  "description, d",
			Usage: "pull-request description",
		},
	},
}

func runCreatePullRequest(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)

/*
  Head      string   `json:"head" binding:"Required"`
  Base      string   `json:"base" binding:"Required"`
  Title     string   `json:"title" binding:"Required"`
  Body      string   `json:"body"`
  Assignee  string   `json:"assignee"`
  Assignees []string `json:"assignees"`
  Milestone int64    `json:"milestone"`
  Labels    []int64  `json:"labels"`
  // swagger:strfmt date-time
  Deadline *time.Time `json:"due_date"`
*/

	pr, err := login.Client().CreatePullRequest(owner, repo, gitea.CreatePullRequestOption{
		Head: ctx.String("head"),
    Base: ctx.String("base"),
    Title: ctx.String("title"),
    Body: ctx.String("body"),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("#%d %s\n%s created %s\n\n%s", pr.Index,
		pr.Title,
		pr.Poster.UserName,
		pr.Created.Format("2006-01-02 15:04:05"),
		pr.Body,
	)
	return nil
}

