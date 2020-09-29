// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package login

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/intern"

	"github.com/urfave/cli/v2"
)

// CmdLoginSetDefault represents to login a gitea server.
var CmdLoginSetDefault = cli.Command{
	Name:        "default",
	Usage:       "Get or Set Default Login",
	Description: `Get or Set Default Login`,
	ArgsUsage:   "<Login>",
	Action:      runLoginSetDefault,
	Flags:       []cli.Flag{&flags.OutputFlag},
}

func runLoginSetDefault(ctx *cli.Context) error {
	if err := intern.LoadConfig(); err != nil {
		return err
	}
	if ctx.Args().Len() == 0 {
		l, err := intern.GetDefaultLogin()
		if err != nil {
			return err
		}
		fmt.Printf("Default Login: %s\n", l.Name)
		return nil
	}
	loginExist := false
	for i := range intern.Config.Logins {
		intern.Config.Logins[i].Default = false
		if intern.Config.Logins[i].Name == ctx.Args().First() {
			intern.Config.Logins[i].Default = true
			loginExist = true
		}
	}

	if !loginExist {
		return fmt.Errorf("login '%s' not found", ctx.Args().First())
	}

	return intern.SaveConfig()
}
