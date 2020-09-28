// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/urfave/cli/v2"
)

// CmdLogin represents to login a gitea server.
var CmdLogin = cli.Command{
	Name:        "login",
	Usage:       "Log in to a Gitea server",
	Description: `Log in to a Gitea server`,
	Action:      runLoginAddInteractive, // TODO show list if no arg & detail if login as arg
	Subcommands: []*cli.Command{
		&cmdLoginList,
		&cmdLoginAdd,
		&cmdLoginEdit,
		&cmdLoginSetDefault,
	},
}
