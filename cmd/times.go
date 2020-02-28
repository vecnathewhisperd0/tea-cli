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

	"code.gitea.io/sdk/gitea"

	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

// CmdTrackedTimes represents the command to operate repositories' times.
var CmdTrackedTimes = cli.Command{
	Name:        "times",
	Aliases: []string{"time"},
	Usage:       "Operate on tracked times of a repository's issues & pulls",
	Description: `Operate on tracked times of a repository's issues & pulls.
		 Depending on your permissions on the repository, only your own tracked
		 times might be listed.`,
	ArgsUsage:   "[username | #issue]",
	Action:      runTrackedTimes,
	Subcommands: []*cli.Command{
		&CmdTrackedTimesAdd,
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
	login, owner, repo := initCommand()
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
		times, err = client.GetRepoTrackedTimes(owner, repo)
	} else if strings.HasPrefix(user, "#") {
		// get all tracked times on the specified issue
		issue, err2 := strconv.ParseInt(user[1:], 10, 64)
		if err2 != nil {
			return err2
		}
		times, err = client.ListTrackedTimes(owner, repo, issue)
	} else {
		// get all tracked times by the specified user
		times, err = client.GetUserTrackedTimes(owner, repo, user)
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

	printTrackedTimes(times, outputValue, from, until, ctx.Bool("total"))
	return nil
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
				time.Duration(1e9 * t.Time).String(),
			},
		)
	}

	if printTotal {
		outputValues = append(outputValues, []string{
			"TOTAL", "", "", time.Duration(1e9 * totalDuration).String(),
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
	login, owner, repo := initCommand()

	if ctx.Args().Len() < 2 {
		return fmt.Errorf("No issue or duration specified.\nUsage:\t%s", ctx.Command.UsageText)
	}

	issueStr := ctx.Args().First()
	if strings.HasPrefix(issueStr, "#") {
		issueStr = issueStr[1:]
	}
	issue, err := strconv.ParseInt(issueStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	duration, err := time.ParseDuration(strings.Join(ctx.Args().Tail(), ""))
	if err != nil {
		log.Fatal(err)
	}

	_, err = login.Client().AddTime(owner, repo, issue, gitea.AddTimeOption{
		Time: int64(duration.Seconds()),
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
