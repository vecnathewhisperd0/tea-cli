// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdIssuesClose represents a sub command of issues to close an issue
var CmdIssuesClose = cli.Command{
	Name:        "close",
	Usage:       "Change state of an issue to 'closed'",
	Description: `Change state of an issue to 'closed'`,
	ArgsUsage:   "<issue index>",
	Action: func(ctx *cli.Context) error {
		var s = gitea.StateClosed
		return editIssueState(ctx, gitea.EditIssueOption{State: &s})
	},
	Flags: AllDefaultFlags,
}
