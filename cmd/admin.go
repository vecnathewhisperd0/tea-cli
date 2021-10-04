// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/urfave/cli/v2"
)

// CmdAdmin represents the admin sub-commands
var CmdAdmin = cli.Command{
	Name:     "admin",
	Aliases:  []string{"a"},
	Category: catAdmin,
	Action: func(cmd *cli.Context) error {
		// TODO: this is just a stub for all admin actions
		//       there is no default admin action
		return nil
	},
	Subcommands: []*cli.Command{
		&cmdAdminUsers,
	},
}
