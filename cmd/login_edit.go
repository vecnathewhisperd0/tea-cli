// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/modules/intern"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

// cmdLoginEdit represents to login a gitea server.
var cmdLoginEdit = cli.Command{
	Name:        "edit",
	Usage:       "Edit Gitea logins",
	Description: `Edit Gitea logins`,
	Action:      runLoginEdit,
	Flags:       []cli.Flag{&OutputFlag},
}

func runLoginEdit(ctx *cli.Context) error {
	return open.Start(intern.GetConfigPath())
}
