// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/admin/users"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

var cmdAdminUsers = cli.Command{
	Name:     "user",
	Aliases:  []string{"u"},
	Category: catAdmin,
	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() == 1 {
			return runAdminUserDetail(ctx, ctx.Args().First())
		}
		return users.RunUserList(ctx)
	},
	Subcommands: []*cli.Command{
		&users.CmdUserList,
	},
	Flags: users.CmdUserList.Flags,
}

func runAdminUserDetail(cmd *cli.Context, u string) error {
	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()
	user, _, err := client.GetUserInfo(u)
	if err != nil {
		return err
	}

	print.UserDetails(user)
	return nil
}
