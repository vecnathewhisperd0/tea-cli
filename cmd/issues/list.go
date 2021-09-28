// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"time"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

var issueFieldsFlag = flags.FieldsFlag(print.IssueFields, []string{
	"index", "title", "state", "author", "milestone", "labels",
})

// CmdIssuesList represents a sub command of issues to list issues
var CmdIssuesList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List issues of the repository",
	Description: `List issues of the repository`,
	Action:      RunIssuesList,
	Flags:       append([]cli.Flag{issueFieldsFlag}, flags.IssuePRFlags...),
}

// RunIssuesList list issues
func RunIssuesList(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "open":
		state = gitea.StateOpen
	case "closed":
		state = gitea.StateClosed
	}

	var err error
	var from, until time.Time
	if ctx.IsSet("from") {
		from, err = dateparse.ParseLocal(ctx.String("from"))
		if err != nil {
			return err
		}
	}
	if ctx.IsSet("until") {
		until, err = dateparse.ParseLocal(ctx.String("until"))
		if err != nil {
			return err
		}
	}

	// ignore error, as we don't do any input validation on these flags
	labels, _ := flags.LabelFilterFlag.GetValues(cmd)
	milestones, _ := flags.MilestoneFilterFlag.GetValues(cmd)

	issues, _, err := ctx.Login.Client().ListRepoIssues(ctx.Owner, ctx.Repo, gitea.ListIssueOption{
		ListOptions: ctx.GetListOptions(),
		State:       state,
		Type:        gitea.IssueTypeIssue,
		KeyWord:     ctx.String("keyword"),
		CreatedBy:   ctx.String("author"),
		AssignedBy:  ctx.String("assigned-to"),
		MentionedBy: ctx.String("mentions"),
		Labels:      labels,
		Milestones:  milestones,
		Since:       from,
		Before:      until,
	})

	if err != nil {
		return err
	}

	fields, err := issueFieldsFlag.GetValues(cmd)
	if err != nil {
		return err
	}

	print.IssuesPullsList(issues, ctx.Output, fields)
	return nil
}
