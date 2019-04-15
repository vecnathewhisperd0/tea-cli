// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	//"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "Operate with repositories",
	Description: `Operate with repositories`,
	Action:      runRepos,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login, optional when inside a gitea repository",
		},
		cli.BoolFlag{
			Name:  "noheader",
			Usage: "Skip header on output.",
		},
	},
}

func runRepos(ctx *cli.Context) error {
	login := initCommandLoginOnly(ctx)

	rps, err := login.Client().ListMyRepos()

	if err != nil {
		log.Fatal(err)
	}

	if len(rps) == 0 {
		fmt.Println("No repositories found")
		return nil
	}

	if noheader := getGlobalFlag(ctx, "noheader"); noheader != "true" {
		fmt.Printf("Name\tType/Mode\tSSH-URL\tOwner\n")
	}
	for _, rp := range rps {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", rp.FullName, mode, rp.SSHURL, rp.Owner.UserName)
	}

	return nil
}

func initCommandLoginOnly(ctx *cli.Context) (*Login) {
	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("load config file failed", yamlConfigPath)
	}

	var login *Login
	if loginFlag := getGlobalFlag(ctx, "login"); loginFlag == "" {
		login, err = getActiveLogin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		login = getLoginByName(loginFlag)
		if login == nil {
			log.Fatal("indicated login name", loginFlag, "does not exist")
		}
	}
	return login
}
