// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"

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
	))
}

// IssuesList prints a listing of issues
func IssuesList(issues []*gitea.Issue, output string) {
	// TODO: make fields selectable
	fields := []string{"index", "title", "state", "author", "milestone", "updated"}
	printIssues(issues, output, fields)
}

// IssuesPullsList prints a listing of issues & pulls
func IssuesPullsList(issues []*gitea.Issue, output string) {
	// TODO: make fields selectable
	fields := []string{"index", "state", "kind", "author", "updated", "title"}
	printIssues(issues, output, fields)
}

// IssueFields are all available fields to print with IssuesList()
var IssueFields = []string{
	"index",
	"state",
	"kind",
	"author",
	"updated",
	"title",
	"milestone",
}

func printIssues(issues []*gitea.Issue, output string, fields []string) {
	var printables = make([]printable, len(issues))
	for i, x := range issues {
		printables[i] = &printableIssue{x}
	}

	t := tableFromItems(fields, printables)
	t.print(output)
}

type printableIssue struct{ *gitea.Issue }

func (x printableIssue) FormatField(field string) string {
	switch field {
	case "index":
		return strconv.FormatInt(x.Index, 10)
	case "state":
		return string(x.State)
	case "kind":
		if x.PullRequest != nil {
			return "Pull"
		}
		return "Issue"
	case "author":
		name := x.Poster.FullName
		if len(name) == 0 {
			return x.Poster.UserName
		}
		return name
	case "updated":
		return FormatTime(x.Updated)
	case "title":
		return x.Title
	case "milestone":
		if x.Milestone != nil {
			return x.Milestone.Title
		}
		return ""
	}
	return ""
}
