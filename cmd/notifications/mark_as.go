// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"github.com/urfave/cli/v2"
)

// CmdNotificationsRead represents a sub command of notifications to list read notifications
var CmdNotificationsMarkRead = cli.Command{
	Name:        "read",
	Aliases:     []string{},
	Usage:       "Show read notifications only",
	Description: `Show read notifications only`,
	Action:      RunNotificationsRead,
	Flags:       flags.AllDefaultFlags,
}

// RunNotificationsRead will show notifications with status read.
func RunNotificationsRead(ctx *cli.Context) error {
	var statuses = []gitea.NotifyStatus{gitea.NotifyStatusRead}
	return listNotifications(ctx, statuses, []gitea.NotifySubjectType{})
}

// CmdNotificationsUnread represents a sub command of notifications to list unread notifications.
var CmdNotificationsMarkUnread = cli.Command{
	Name:        "unread",
	Aliases:     []string{},
	Usage:       "Show unread notifications only",
	Description: `Show unread notifications only`,
	Action:      RunNotificationsUnread,
	Flags:       flags.AllDefaultFlags,
}

// RunNotificationsUnread will show notifications with status unread.
func RunNotificationsUnread(ctx *cli.Context) error {
	var statuses = []gitea.NotifyStatus{gitea.NotifyStatusUnread}
	return listNotifications(ctx, statuses, []gitea.NotifySubjectType{})
}
