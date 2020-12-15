// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"
	"strings"

	"code.gitea.io/sdk/gitea"
)

var ciStatusSymbols = map[gitea.StatusState]string{
	gitea.StatusSuccess: "✓ ",
	gitea.StatusPending: "⭮ ",
	gitea.StatusWarning: "⚠ ",
	gitea.StatusError:   "✘ ",
	gitea.StatusFailure: "❌ ",
}

// PullDetails print an pull rendered to stdout
func PullDetails(pr *gitea.PullRequest, reviews []*gitea.PullReview, ciStatus *gitea.CombinedStatus) {
	base := pr.Base.Name
	head := pr.Head.Name
	if pr.Head.RepoID != pr.Base.RepoID {
		if pr.Head.Repository != nil {
			head = pr.Head.Repository.Owner.UserName + ":" + head
		} else {
			head = "delete:" + head
		}
	}

	out := fmt.Sprintf(
		"# #%d %s (%s)\n@%s created %s\t**%s** <- **%s**\n\n%s\n\n",
		pr.Index,
		pr.Title,
		pr.State,
		pr.Poster.UserName,
		FormatTime(*pr.Created),
		base,
		head,
		pr.Body,
	)

	if len(reviews) != 0 {
		out += "---\n"
		revMap := make(map[gitea.ReviewStateType][]string)
		for _, review := range reviews {
			switch review.State {
			case gitea.ReviewStateApproved,
				gitea.ReviewStateRequestChanges,
				gitea.ReviewStateRequestReview:
				u := revMap[review.State]
				revMap[review.State] = append(u, review.Reviewer.UserName)
			}
		}
		for k, v := range revMap {
			out += fmt.Sprintf("- %s by @%s\n", k, strings.Join(v, ", @"))
		}
	}

	if ciStatus != nil {
		var summary, errors string
		for _, s := range ciStatus.Statuses {
			summary += ciStatusSymbols[s.State]
			if s.State != gitea.StatusSuccess {
				errors += fmt.Sprintf("  - [**%s**:\t%s](%s)\n", s.Context, s.Description, s.TargetURL)
			}
		}
		if len(ciStatus.Statuses) != 0 {
			out += fmt.Sprintf("- CI: %s\n%s", summary, errors)
		}
	}

	if pr.State == gitea.StateOpen {
		if pr.Mergeable {
			out += "- No Conflicts\n"
		} else {
			out += "- **Conflicting files**\n"
		}
	}

	outputMarkdown(out)
}

// PullsList prints a listing of pulls
func PullsList(prs []*gitea.PullRequest, output string) {
	t := tableWithHeader(
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	)

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
		t.addRow(
			strconv.FormatInt(pr.Index, 10),
			pr.Title,
			string(pr.State),
			author,
			mile,
			FormatTime(*pr.Updated),
		)
	}

	t.print(output)
}
