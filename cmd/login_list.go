// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/tea/modules/intern"

	"github.com/urfave/cli/v2"
)

// CmdLogin represents to login a gitea server.
var cmdLoginList = cli.Command{
	Name:        "ls",
	Usage:       "List Gitea logins",
	Description: `List Gitea logins`,
	Action:      runLoginList,
	Flags:       []cli.Flag{&OutputFlag},
}

func runLoginList(ctx *cli.Context) error {
	err := intern.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Name",
		"URL",
		"SSHHost",
		"User",
		"Default",
	}

	var values [][]string

	for _, l := range intern.Config.Logins {
		values = append(values, []string{
			l.Name,
			l.URL,
			l.GetSSHHost(),
			l.User,
			fmt.Sprint(l.Default),
		})
	}

	Output(globalOutputValue, headers, values)

	return nil
}
