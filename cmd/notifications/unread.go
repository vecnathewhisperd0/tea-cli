// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"github.com/urfave/cli/v2"
)

// CmdNotificationsUnread represents a sub command of notifications to list unread notifications.
var CmdNotificationsUnread = cli.Command{
	Name:        "unread",
	Aliases:     []string{},
	Usage:       "show unread notifications",
	Description: `show unread notifications`,
	Action:      RunNotificationsUnread,
	Flags:       flags.NotificationFlags,
}

// RunNotificationsUnread will show notifications with status unread.
func RunNotificationsUnread(ctx *cli.Context) error {
	var statuses = []gitea.NotifyStatus{gitea.NotifyStatusUnread}
	return listNotifications(ctx, statuses)
}
