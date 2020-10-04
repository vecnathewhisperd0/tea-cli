// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/charmbracelet/glamour"
)

// PullDetails print an pull rendered to stdout
func PullDetails(pr *gitea.PullRequest) {

	in := fmt.Sprintf("# #%d %s (%s)\n%s created %s\n\n%s\n", pr.Index,
		pr.Title,
		pr.State,
		pr.Poster.UserName,
		FormatTime(*pr.Created),
		pr.Body,
	)
	out, err := glamour.Render(in, getGlamourTheme())
	if err != nil {
		// TODO: better Error handling
		fmt.Printf("Error:\n%v\n\n", err)
		return
	}
	fmt.Print(out)
}
