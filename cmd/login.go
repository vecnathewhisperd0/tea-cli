// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"

	"code.gitea.io/tea/cmd/login"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

// CmdLogin represents to login a gitea server.
var CmdLogin = cli.Command{
	Name:        "logins",
	Aliases:     []string{"login"},
	Category:    catSetup,
	Usage:       "Log in to a Gitea server",
	Description: `Log in to a Gitea server`,
	ArgsUsage:   "[<login name>]",
	Action:      runLogins,
	Subcommands: []*cli.Command{
		&login.CmdLoginList,
		&login.CmdLoginAdd,
		&login.CmdLoginEdit,
		&login.CmdLoginDelete,
		&login.CmdLoginSetDefault,
		&login.CmdLoginHelper,
	},
}

func runLogins(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runLoginDetail(ctx.Args().First())
	}
	return login.RunLoginList(ctx)
}

func runLoginDetail(name string) error {
	l := config.GetLoginByName(name)
	if l == nil {
		fmt.Printf("Login '%s' do not exist\n\n", name)
		return nil
	}

	print.LoginDetails(l)
	return nil
}
