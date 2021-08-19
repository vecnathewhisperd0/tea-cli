// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"log"

	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// listNotifications will get the notifications based on status
func listNotifications(cmd *cli.Context, status []gitea.NotifyStatus) error {
	var news []*gitea.NotificationThread
	var err error

	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()
	all := ctx.Bool("for-user")

	// This enforces pagination (see https://github.com/go-gitea/gitea/issues/16733)
	listOpts := ctx.GetListOptions()
	if listOpts.Page == 0 {
		listOpts.Page = 1
	}

	if all {
		news, _, err = client.ListNotifications(gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
			// TODO: SubjectTypes
		})
	} else {
		ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
		news, _, err = client.ListRepoNotifications(ctx.Owner, ctx.Repo, gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	}
	if err != nil {
		log.Fatal(err)
	}

	print.NotificationsList(news, ctx.Output, all)
	return nil
}
