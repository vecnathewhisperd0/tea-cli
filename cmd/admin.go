// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/urfave/cli/v2"
)

// CmdAdmin represents the namespace of admin commands.
// The command itself has no functionality, but hosts subcommands.
var CmdAdmin = cli.Command{
	Name:     "admin",
	Aliases:  []string{"a"},
	Category: catHelpers,
	Action: func(cmd *cli.Context) error {
		return cli.ShowSubcommandHelp(cmd)
	},
	Subcommands: []*cli.Command{
		&cmdAdminUsers,
	},
}
