// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"log"
	"os"

	"code.gitea.io/tea/modules/intern"

	"github.com/urfave/cli/v2"
)

// CmdLogout represents to logout a gitea server.
var CmdLogout = cli.Command{
	Name:        "logout",
	Usage:       "Log out from a Gitea server",
	Description: `Log out from a Gitea server`,
	Action:      runLogout,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Login name to remove",
		},
	},
}

func runLogout(ctx *cli.Context) error {
	var name string
	if len(os.Args) == 3 {
		name = os.Args[2]
	} else if ctx.IsSet("name") {
		name = ctx.String("name")
	} else {
		return errors.New("Please specify a login name")
	}

	err := intern.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var idx = -1
	for i, l := range intern.Config.Logins {
		if l.Name == name {
			idx = i
			break
		}
	}
	if idx > -1 {
		intern.Config.Logins = append(intern.Config.Logins[:idx], intern.Config.Logins[idx+1:]...)
		err = intern.SaveConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
