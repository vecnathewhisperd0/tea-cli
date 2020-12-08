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

// CmdNotificationsList represents a sub command of notifications to list notifications
var CmdNotificationsList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List notifications....",
	Description: `List notification....`,
	Action:      RunNotificationsList,
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
func RunNotificationsList(ctx *cli.Context) error {
	return task.ListNotifications(ctx, []gitea.NotifyStatus{})
}
