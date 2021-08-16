// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdNotificationsList represents a sub command of notifications to list notifications
var CmdNotificationsList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List notifications",
	Description: `List notifications`,
	Action:      RunNotificationsList,
	//	Flags:       flags.NotificationFlags,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "show all notifications of related gitea instance",
		},
		&cli.StringFlag{
			Name:        "state",
			Usage:       "set notification state (default is all), pinned,read,unread",
			DefaultText: "xxx",
		},
	}),
}

// notif ls
// notif ls --state all
// notif ls --state pinned
// notif ls --state read
// notif ls --state unread

// RunNotificationsList list notifications
func RunNotificationsList(ctx *cli.Context) error {
	// --states all

	states := []gitea.NotifyStatus{}

	switch ctx.String("state") {
	case "all":
		states = []gitea.NotifyStatus{gitea.NotifyStatusPinned, gitea.NotifyStatusRead, gitea.NotifyStatusUnread}
	case "pinned":
		states = []gitea.NotifyStatus{gitea.NotifyStatusPinned}
	case "read":
		states = []gitea.NotifyStatus{gitea.NotifyStatusRead}
	case "unread":
		states = []gitea.NotifyStatus{gitea.NotifyStatusUnread}
	default:
		return fmt.Errorf("invalid notification state type '%s'. valid: all, pinned,read,unread", ctx.String("state"))
	}

	return listNotifications(ctx, states)
}
