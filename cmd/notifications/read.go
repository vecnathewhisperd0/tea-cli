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

// CmdNotificationsRead represents a sub command of notifications to list read notifications
var CmdNotificationsRead = cli.Command{
	Name:        "read",
	Aliases:     []string{},
	Usage:       "show read notifications instead",
	Description: `show read notifications instead`,
	Action:      RunNotificationsRead,
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

// RunNotificationsRead will show notifications with status read.
func RunNotificationsRead(ctx *cli.Context) error {
	var statuses = []gitea.NotifyStatus{gitea.NotifyStatusRead}
	return task.ListNotifications(ctx, statuses)
}
