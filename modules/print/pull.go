// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"

	"code.gitea.io/tea/cmd/flags"

	"code.gitea.io/sdk/gitea"
)

// PullDetails print an pull rendered to stdout
func PullDetails(pr *gitea.PullRequest) {
	OutputMarkdown(fmt.Sprintf(
		"# #%d %s (%s)\n%s created %s\n\n%s\n",
		pr.Index,
		pr.Title,
		pr.State,
		pr.Poster.UserName,
		FormatTime(*pr.Created),
		pr.Body,
	))
}

// PullsList prints a listing of pulls
func PullsList(prs []*gitea.PullRequest) {
	var values [][]string
	headers := []string{
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	}

	if len(prs) == 0 {
		OutputList(flags.GlobalOutputValue, headers, values)
		return
	}

	for _, pr := range prs {
		if pr == nil {
			continue
		}
		author := pr.Poster.FullName
		if len(author) == 0 {
			author = pr.Poster.UserName
		}
		mile := ""
		if pr.Milestone != nil {
			mile = pr.Milestone.Title
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(pr.Index, 10),
				pr.Title,
				string(pr.State),
				author,
				mile,
				FormatTime(*pr.Updated),
			},
		)
	}

	OutputList(flags.GlobalOutputValue, headers, values)
}
