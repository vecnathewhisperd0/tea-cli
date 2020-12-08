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

// CmdNotificationsUnread represents a sub command of notifications to list unread notifications.
var CmdNotificationsUnread = cli.Command{
	Name:        "unread",
	Aliases:     []string{},
	Usage:       "show unread notifications",
	Description: `show unread notifications`,
	Action:      RunNotificationsUnread,
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

// RunNotificationsList list notifications
func RunNotificationsUnread(ctx *cli.Context) error {
	var statuses = []gitea.NotifyStatus{gitea.NotifyStatusUnread}
	return task.ListNotifications(ctx, statuses)
}
