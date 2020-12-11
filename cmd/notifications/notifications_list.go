// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"fmt"
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

//listNotifications will get the notifications based on status
func listNotifications(ctx *cli.Context, status []gitea.NotifyStatus) error {

	//This enforces pagination.
	listOpts := flags.GetListOptions(ctx)
	if listOpts.Page == 0 {
		listOpts.Page = 1
	}

	var news []*gitea.NotificationThread
	var err error

	var allRelated = ctx.Bool("all")
	fmt.Printf("allRelated: %t\n", allRelated)

	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	if allRelated {
		fmt.Printf("login: %s owner: %s repo:%s\n", login.Name, owner, repo)
		news, _, err = login.Client().ListNotifications(gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	} else {
		fmt.Printf("login: %s owner: %s repo:%s\n", login.Name, owner, repo)
		news, _, err = login.Client().ListRepoNotifications(owner, repo, gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	}
	if err != nil {
		log.Fatal(err)
	}

	print.NotificationsList(news, flags.GlobalOutputValue, allRelated)
	return nil
}
