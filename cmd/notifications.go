// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/tea/cmd/notifications"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdNotifications is the main command to operate with notifications
var CmdNotifications = cli.Command{
	Name:        "notifications",
	Aliases:     []string{"notification", "n"},
	Usage:       "Show notifications",
	Description: "Show notifications, by default based of the current repo",
	Action:      runNotifications,
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
			Name:        "state",
			Usage:       "set milestone state (default is open)",
			DefaultText: "open",
		},
	}),
}

func runNotifications(cmd *cli.Context) error {
	var news []*gitea.NotificationThread
	var err error

	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()

	listOpts := ctx.GetListOptions()
	if listOpts.Page == 0 {
		listOpts.Page = 1
	}

	var status []gitea.NotifyStatus
	if ctx.Bool("read") {
		status = []gitea.NotifyStatus{gitea.NotifyStatusRead}
	}
	if ctx.Bool("pinned") {
		status = append(status, gitea.NotifyStatusPinned)
	}

	if ctx.Bool("all") {
		news, _, err = client.ListNotifications(gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	} else {
		ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
		news, _, err = client.ListRepoNotifications(ctx.Owner, ctx.Repo, gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	}
	if err != nil {
		return err
	}

	print.NotificationsList(news, ctx.Output, ctx.Bool("all"))
	return nil
}

func runNotificationsDetails(ctx *cli.Context) error {
	log.Fatal("Use tea notif --help to see all available commands.")
	return nil
}
