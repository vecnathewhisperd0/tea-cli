// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"github.com/urfave/cli/v2"
)

// CmdNotificationsUnpin represents a sub command of notifications to unpin a notification
var CmdNotificationsUnpin = cli.Command{
	Name:        "unpin",
	Usage:       "Remove a notification from pins",
	Description: `Remove a notification from pins`,
	Action:      RunNotificationsPinned,
	Flags:       flags.AllDefaultFlags,
}

// CmdNotificationsPin represents a sub command of notifications to pin a notification
var CmdNotificationsPin = cli.Command{
	Name:        "pin",
	Aliases:     []string{"p"},
	Usage:       "Save a notification as pin",
	Description: `Save a notification as pin`,
	Action:      RunNotificationsPinned,
	Flags:       flags.AllDefaultFlags,
}

// RunNotificationsPinned will show notifications with status pinned.
func RunNotificationsPinned(ctx *cli.Context) error {
	var statuses = []gitea.NotifyStatus{gitea.NotifyStatusPinned}
	return listNotifications(ctx, statuses, []gitea.NotifySubjectType{})
}
