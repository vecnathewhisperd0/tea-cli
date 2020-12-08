// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/task"
	"github.com/urfave/cli/v2"
)

// CmdNotificationsPinned represents a sub command of notifications to list pinned notifications
var CmdNotificationsPinned = cli.Command{
	Name:        "pinned",
	Aliases:     []string{"pin"},
	Usage:       "show pinned notifications",
	Description: `show pinned notifications`,
	Action:      RunNotificationsPinned,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "show all notifications of related gitea instance",
		},
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by milestone state (all|open|closed)",
			DefaultText: "open",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunNotificationsPinned will show notifications with status pinned.
func RunNotificationsPinned(ctx *cli.Context) error {
	var statuses = []gitea.NotifyStatus{gitea.NotifyStatusPinned}
	return task.ListNotifications(ctx, statuses)
}
