// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/issues"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

var prFieldsFlag = flags.FieldsFlag(print.IssueFields, []string{
	"index", "title", "state", "author", "milestone", "updated",
})

// CmdPullsList represents a sub command of issues to list pulls
var CmdPullsList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List pull requests of the repository",
	Description: `List pull requests of the repository`,
	Action:      RunPullsList,
	Flags:       append([]cli.Flag{prFieldsFlag}, flags.IssuePRFlags...),
}

// RunPullsList return list of pulls
func RunPullsList(cmd *cli.Context) error {
	return issues.DoIssuePRListing(cmd, gitea.IssueTypePull, prFieldsFlag)
}
