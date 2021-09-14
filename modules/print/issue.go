// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
)

// IssueDetails print an issue rendered to stdout
func IssueDetails(issue *gitea.Issue) {
	outputMarkdown(fmt.Sprintf(
		"# #%d %s (%s)\n@%s created %s\n\n%s\n",
		issue.Index,
		issue.Title,
		issue.State,
		issue.Poster.UserName,
		FormatTime(issue.Created),
		issue.Body,
	), getRepoURL(issue.HTMLURL))
}

// IssuesPullsList prints a listing of issues & pulls
func IssuesPullsList(issues []*gitea.Issue, output string, fields []string) {
	printIssues(issues, output, fields)
}

// IssueFields are all available fields to print with IssuesList()
var IssueFields = []string{
	"index",
	"state",
	"kind",
	"author",
	"author-id",
	"url",

	"title",
	"body",

	"created",
	"updated",
	"deadline",

	"assignees",
	"milestone",
	"labels",
	"comments",
}

func printIssues(issues []*gitea.Issue, output string, fields []string) {
	labelMap := map[int64]string{}
	var printables = make([]printable, len(issues))

	for i, x := range issues {
		// pre-serialize labels for performance
		for _, label := range x.Labels {
			if _, ok := labelMap[label.ID]; !ok {
				labelMap[label.ID] = formatLabel(label, !isMachineReadable(output), "")
			}
		}
		// store items with printable interface
		printables[i] = &printableIssue{x, &labelMap}
	}

	t := tableFromItems(fields, printables)
	t.print(output)
}

type printableIssue struct {
	*gitea.Issue
	formattedLabels *map[int64]string
}

func (x printableIssue) FormatField(field string) string {
	switch field {
	case "index":
		return fmt.Sprintf("%d", x.Index)
	case "state":
		return string(x.State)
	case "kind":
		if x.PullRequest != nil {
			return "Pull"
		}
		return "Issue"
	case "author":
		return formatUserName(x.Poster)
	case "author-id":
		return x.Poster.UserName
	case "url":
		return x.HTMLURL
	case "title":
		return x.Title
	case "body":
		return x.Body
	case "created":
		return FormatTime(x.Created)
	case "updated":
		return FormatTime(x.Updated)
	case "deadline":
		if x.Deadline == nil {
			return ""
		}
		return FormatTime(*x.Deadline)
	case "milestone":
		if x.Milestone != nil {
			return x.Milestone.Title
		}
		return ""
	case "labels":
		var labels = make([]string, len(x.Labels))
		for i, l := range x.Labels {
			labels[i] = (*x.formattedLabels)[l.ID]
		}
		return strings.Join(labels, " ")
	case "assignees":
		var assignees = make([]string, len(x.Assignees))
		for i, a := range x.Assignees {
			assignees[i] = formatUserName(a)
		}
		return strings.Join(assignees, " ")
	case "comments":
		return fmt.Sprintf("%d", x.Comments)
	}
	return ""
}
