// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/urfave/cli/v2"
)

// CmdMilestonesReopen represents a sub command of milestones to open an milestone
var CmdMilestonesReopen = cli.Command{
	Name:        "reopen",
	Aliases:     []string{"open"},
	Usage:       "Change state of an milestone to 'open'",
	Description: `Change state of an milestone to 'open'`,
	ArgsUsage:   "<milestone name>",
	Action: func(ctx *cli.Context) error {
		return editMilestoneStatus(ctx, false)
	},
	Flags: AllDefaultFlags,
}
