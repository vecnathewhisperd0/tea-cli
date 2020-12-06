// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
)

// IssueDetails print an issue rendered to stdout
func IssueDetails(issue *gitea.Issue) {
	OutputMarkdown(fmt.Sprintf(
		"# #%d %s (%s)\n%s created %s\n\n%s\n",
		issue.Index,
		issue.Title,
		issue.State,
		issue.Poster.UserName,
		FormatTime(issue.Created),
		issue.Body,
	))
}

// IssuesList prints a listing of issues
func IssuesList(issues []*gitea.Issue) {
	var values [][]string
	headers := []string{
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	}

	if len(issues) == 0 {
		OutputList(flags.GlobalOutputValue, headers, values)
		return
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
				FormatTime(issue.Updated),
			},
		)
	}
	OutputList(flags.GlobalOutputValue, headers, values)
}
