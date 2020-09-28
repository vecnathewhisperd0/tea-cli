// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"code.gitea.io/tea/modules/intern"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdTrackedTimesAdd represents a sub command of times to add time to an issue
var CmdTrackedTimesAdd = cli.Command{
	Name:      "add",
	Usage:     "Track spent time on an issue",
	UsageText: "tea times add <issue> <duration>",
	Description: `Track spent time on an issue
	 Example:
		tea times add 1 1h25m
	`,
	Action: runTrackedTimesAdd,
	Flags:  LoginRepoFlags,
}

func runTrackedTimesAdd(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	if ctx.Args().Len() < 2 {
		return fmt.Errorf("No issue or duration specified.\nUsage:\t%s", ctx.Command.UsageText)
	}

	issue, err := argToIndex(ctx.Args().First())
	if err != nil {
		log.Fatal(err)
	}

	duration, err := time.ParseDuration(strings.Join(ctx.Args().Tail(), ""))
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = login.Client().AddTime(owner, repo, issue, gitea.AddTimeOption{
		Time: int64(duration.Seconds()),
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
