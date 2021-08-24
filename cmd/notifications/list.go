// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"fmt"
	"os"

	"code.gitea.io/tea/cmd/flags"

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
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "mine",
			Aliases: []string{"m"},
			Usage:   "Show notifications across all your repositories instead of the current repository only",
		},
		&cli.StringFlag{
			Name:        "state",
			Aliases:     []string{"s"},
			Usage:       "filter by notification state (pinned,read,unread,all)",
			DefaultText: "pinned + unread",
		},
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "filter by subject type (repo,issue,pr,commit)",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunNotificationsList list notifications
func RunNotificationsList(ctx *cli.Context) error {
	var states []gitea.NotifyStatus
	var types []gitea.NotifySubjectType

	// FIXME: make these commaseparated slice flags..?

	switch ctx.String("state") {
	case "":
		states = []gitea.NotifyStatus{}
	case "unread":
		states = []gitea.NotifyStatus{gitea.NotifyStatusUnread}
	case "pinned":
		states = []gitea.NotifyStatus{gitea.NotifyStatusPinned}
	case "read":
		states = []gitea.NotifyStatus{gitea.NotifyStatusRead}
	case "all":
		states = []gitea.NotifyStatus{gitea.NotifyStatusPinned, gitea.NotifyStatusRead, gitea.NotifyStatusUnread}
	default:
		fmt.Printf("Unknown state '%s'\n", ctx.String("state"))
		os.Exit(1)
	}

	switch ctx.String("type") {
	case "":
		types = []gitea.NotifySubjectType{}
	case "commit":
		types = []gitea.NotifySubjectType{gitea.NotifySubjectCommit}
	case "issue":
		types = []gitea.NotifySubjectType{gitea.NotifySubjectIssue}
	case "pull", "pr":
		types = []gitea.NotifySubjectType{gitea.NotifySubjectPull}
	case "repo":
		types = []gitea.NotifySubjectType{gitea.NotifySubjectRepository}
	default:
		fmt.Printf("Unknown type '%s'\n", ctx.String("type"))
		os.Exit(1)
	}

	return listNotifications(ctx, states, types)
}
