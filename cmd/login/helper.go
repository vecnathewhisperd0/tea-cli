// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package login

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"
	"github.com/urfave/cli/v2"
)

// CmdLoginHelper represents to login a gitea helper.
var CmdLoginHelper = cli.Command{
	Name:        "helper",
	Aliases:     []string{"git-credential"},
	Usage:       "Git helper",
	Description: `Git helper`,
	Hidden:      true,
	Subcommands: []*cli.Command{
		{
			Name:        "store",
			Description: "Command drops",
			Aliases:     []string{"erase"},
			Action: func(ctx *cli.Context) error {
				return nil
			},
		},
		{
			Name:        "setup",
			Description: "Setup helper to tea authenticate",
			Action: func(ctx *cli.Context) error {
				logins, err := config.GetLogins()
				if err != nil {
					return err
				}
				for loginIndex := range logins {
					login := logins[loginIndex]
					err = task.SetupHelper(login)
					if err != nil {
						return err
					}
					fmt.Printf("Added %s\n", login.Name)
				}
				return nil
			},
		},
		{
			Name:        "get",
			Description: "Get token to auth",
			Action: func(cmd *cli.Context) error {
				wants := map[string]string{}
				s := bufio.NewScanner(os.Stdin)
				for s.Scan() {
					line := s.Text()
					if line == "" {
						break
					}
					parts := strings.SplitN(line, "=", 2)
					if len(parts) < 2 {
						continue
					}
					key, value := parts[0], parts[1]
					if key == "url" {
						u, err := url.Parse(value)
						if err != nil {
							return err
						}
						wants["protocol"] = u.Scheme
						wants["host"] = u.Host
						wants["path"] = u.Path
						wants["username"] = u.User.Username()
						wants["password"], _ = u.User.Password()
					} else {
						wants[key] = value
					}
				}

				if len(wants["host"]) == 0 {
					log.Fatal("Require hostname")
				} else if len(wants["protocol"]) == 0 {
					wants["protocol"] = "http"
				}

				userConfig := config.GetLoginByHost(wants["host"])
				if len(userConfig.Token) == 0 {
					log.Fatal("User no set")
				}

				host, err := url.Parse(userConfig.URL)
				if err != nil {
					return err
				}

				_, err = fmt.Fprintf(os.Stdout, "protocol=%s\nhost=%s\nusername=%s\npassword=%s\n", host.Scheme, host.Host, userConfig.User, userConfig.Token)
				if err != nil {
					return err
				}

				return nil
			},
		},
	},
}
