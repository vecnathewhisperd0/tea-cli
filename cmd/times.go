// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"code.gitea.io/tea/modules/intern"

	"code.gitea.io/sdk/gitea"
	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

// CmdTrackedTimes represents the command to operate repositories' times.
var CmdTrackedTimes = cli.Command{
	Name:    "times",
	Aliases: []string{"time"},
	Usage:   "Operate on tracked times of a repository's issues & pulls",
	Description: `Operate on tracked times of a repository's issues & pulls.
		 Depending on your permissions on the repository, only your own tracked
		 times might be listed.`,
	ArgsUsage: "[username | #issue]",
	Action:    runTrackedTimes,
	Subcommands: []*cli.Command{
		&CmdTrackedTimesAdd,
		&CmdTrackedTimesDelete,
		&CmdTrackedTimesReset,
	},
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "from",
			Aliases: []string{"f"},
			Usage:   "Show only times tracked after this date",
		},
		&cli.StringFlag{
			Name:    "until",
			Aliases: []string{"u"},
			Usage:   "Show only times tracked before this date",
		},
		&cli.BoolFlag{
			Name:    "total",
			Aliases: []string{"t"},
			Usage:   "Print the total duration at the end",
		},
	}, AllDefaultFlags...),
}

func runTrackedTimes(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)
	client := login.Client()

	if err := client.CheckServerVersionConstraint(">= 1.11"); err != nil {
		return err
	}

	var times []*gitea.TrackedTime
	var err error

	user := ctx.Args().First()
	fmt.Println(ctx.Command.ArgsUsage)
	if user == "" {
		// get all tracked times on the repo
		times, _, err = client.GetRepoTrackedTimes(owner, repo)
	} else if strings.HasPrefix(user, "#") {
		// get all tracked times on the specified issue
		issue, err := argToIndex(user)
		if err != nil {
			return err
		}
		times, _, err = client.ListTrackedTimes(owner, repo, issue, gitea.ListTrackedTimesOptions{})
	} else {
		// get all tracked times by the specified user
		times, _, err = client.GetUserTrackedTimes(owner, repo, user)
	}

	if err != nil {
		return err
	}

	var from, until time.Time
	if ctx.String("from") != "" {
		from, err = dateparse.ParseLocal(ctx.String("from"))
		if err != nil {
			return err
		}
	}
	if ctx.String("until") != "" {
		until, err = dateparse.ParseLocal(ctx.String("until"))
		if err != nil {
			return err
		}
	}

	printTrackedTimes(times, globalOutputValue, from, until, ctx.Bool("total"))
	return nil
}

func formatDuration(seconds int64, outputType string) string {
	switch outputType {
	case "yaml":
	case "csv":
		return fmt.Sprint(seconds)
	}
	return time.Duration(1e9 * seconds).String()
}

func printTrackedTimes(times []*gitea.TrackedTime, outputType string, from, until time.Time, printTotal bool) {
	var outputValues [][]string
	var totalDuration int64

	localLoc, err := time.LoadLocation("Local") // local timezone for time formatting
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range times {
		if !from.IsZero() && from.After(t.Created) {
			continue
		}
		if !until.IsZero() && until.Before(t.Created) {
			continue
		}

		totalDuration += t.Time

		outputValues = append(
			outputValues,
			[]string{
				t.Created.In(localLoc).Format("2006-01-02 15:04:05"),
				"#" + strconv.FormatInt(t.Issue.Index, 10),
				t.UserName,
				formatDuration(t.Time, outputType),
			},
		)
	}

	if printTotal {
		outputValues = append(outputValues, []string{
			"TOTAL", "", "", formatDuration(totalDuration, outputType),
		})
	}

	headers := []string{
		"Created",
		"Issue",
		"User",
		"Duration",
	}
	Output(outputType, headers, outputValues)
}
