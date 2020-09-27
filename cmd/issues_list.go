// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"strconv"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/intern"
	"github.com/urfave/cli/v2"
)

// CmdIssuesList represents a sub command of issues to list issues
var CmdIssuesList = cli.Command{
	Name:        "ls",
	Usage:       "List issues of the repository",
	Description: `List issues of the repository`,
	Action:      runIssuesList,
	Flags:       IssuePRFlags,
}

func runIssuesList(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "open":
		state = gitea.StateOpen
	case "closed":
		state = gitea.StateClosed
	}

	issues, _, err := login.Client().ListRepoIssues(owner, repo, gitea.ListIssueOption{
		ListOptions: getListOptions(ctx),
		State:       state,
		Type:        gitea.IssueTypeIssue,
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	}

	var values [][]string

	if len(issues) == 0 {
		Output(globalOutputValue, headers, values)
		return nil
	}

	for _, issue := range issues {
		author := issue.Poster.FullName
		if len(author) == 0 {
			author = issue.Poster.UserName
		}
		mile := ""
		if issue.Milestone != nil {
			mile = issue.Milestone.Title
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(issue.Index, 10),
				issue.Title,
				string(issue.State),
				author,
				mile,
				issue.Updated.Format("2006-01-02 15:04:05"),
			},
		)
	}
	Output(globalOutputValue, headers, values)

	return nil
}
