// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/notifications"

	"github.com/urfave/cli/v2"
)

// CmdNotifications is the main command to operate with notifications
var CmdNotifications = cli.Command{
	Name:        "notifications",
	Aliases:     []string{"notification", "n"},
	Category:    catHelpers,
	Usage:       "Show notifications",
	Description: "Show notifications, by default based of the current repo",
	Action:      notifications.RunNotificationsList,
	Subcommands: []*cli.Command{
		&notifications.CmdNotificationsList,
		&notifications.CmdNotificationsPinned,
		&notifications.CmdNotificationsRead,
		&notifications.CmdNotificationsUnread,
	},
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "show all notifications of related gitea instance",
		},
		&cli.StringFlag{
			Name:  "state",
			Usage: "set notification state (default is all), pinned,read,unread",
		},
	}),
}
