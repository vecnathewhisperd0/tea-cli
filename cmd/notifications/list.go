// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"code.gitea.io/tea/cmd/flags"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

var notifStateFlag = flags.NewCsvFlag("states", "notification states to filter by", []string{"s"},
	[]string{"pinned", "unread", "read"}, []string{"pinned", "unread"})

var notifTypeFlag = flags.NewCsvFlag("types", "subject types to filter by", []string{"t"},
	[]string{"issue", "pull", "repository", "commit"}, nil)

// CmdNotificationsList represents a sub command of notifications to list notifications
var CmdNotificationsList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List notifications",
	Description: `List notifications`,
	Action:      RunNotificationsList,
	Flags: append([]cli.Flag{
		notifStateFlag,
		notifTypeFlag,
		&cli.BoolFlag{
			Name:    "mine",
			Aliases: []string{"m"},
			Usage:   "Show notifications across all your repositories instead of the current repository only",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunNotificationsList list notifications
func RunNotificationsList(ctx *cli.Context) error {
	var states []gitea.NotifyStatus
	statesStr, err := notifStateFlag.GetValues(ctx)
	if err != nil {
		return err
	}
	for _, s := range statesStr {
		states = append(states, gitea.NotifyStatus(s))
	}

	var types []gitea.NotifySubjectType
	typesStr, err := notifTypeFlag.GetValues(ctx)
	if err != nil {
		return err
	}
	for _, t := range typesStr {
		types = append(types, gitea.NotifySubjectType(t))
	}

	return listNotifications(ctx, states, types)
}
