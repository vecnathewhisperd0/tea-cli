// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/charmbracelet/glamour"
)

func IssueDetails(issue *gitea.Issue) {

	in := fmt.Sprintf("# #%d %s (%s)\n%s created %s\n\n%s\n", issue.Index,
		issue.Title,
		issue.State,
		issue.Poster.UserName,
		issue.Created.Format("2006-01-02 15:04:05"),
		issue.Body,
	)
	out, err := glamour.Render(in, getGlamourTheme())
	if err != nil {
		// TODO: better Error handling
		fmt.Printf("Error:\n%v\n\n", err)
		return
	}
	fmt.Print(out)
}
