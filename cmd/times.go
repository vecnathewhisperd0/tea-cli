// Copyright 2019 The Gitea Authors. All rights reserved.
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
	"github.com/urfave/cli"
)

// CmdTrackedTimes represents the command to operate repositories' times.
var CmdTrackedTimes = cli.Command{
	Name:        "times",
	Usage:       "Operate on tracked times of the repository's issues",
	Description: `Operate on tracked times of the repository's issues`,
	ArgsUsage:   "[username | #issue]",
	Action:      runTrackedTimes,
	Subcommands: []cli.Command{
		CmdTrackedTimesAdd,
	},
	Flags: append([]cli.Flag{
		// TODO: add --from --to filters on t.Created
		cli.StringFlag{
			Name:  "from, f",
			Usage: "Show only times tracked after this date",
		},
		cli.StringFlag{
			Name:  "until, u",
			Usage: "Show only times tracked before this date",
		},
		cli.BoolFlag{
			Name:  "total, t",
			Usage: "Print the total duration at the end",
		},
	}, AllDefaultFlags...),
}

func runTrackedTimes(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	var totalDuration int64
	var times []*gitea.TrackedTime
	var err error
	var outputValues [][]string
	headers := []string{
		"Index",
		"Created",
		"Issue", // FIXME: this is the internal issue ID, not the one of the repo....
		"User",  // FIXME: we should print a username!
		"Duration",
	}

	user := ctx.Args().First()
	fmt.Println(ctx.Command.ArgsUsage)
	if user == "" {
		// get all tracked times on the repo
		times, err = login.Client().GetRepoTrackedTimes(owner, repo)
	} else if strings.HasPrefix(user, "#") {
		// get all tracked times on the specified issue
		issue, err2 := strconv.ParseInt(user[1:], 10, 64)
		if err2 != nil {
			log.Fatal(err2)
		}
		times, err = login.Client().ListTrackedTimes(owner, repo, issue)
	} else {
		// get all tracked times by the specified user
		times, err = login.Client().GetUserTrackedTimes(owner, repo, user)
	}

	if err != nil {
		log.Fatal(err)
	}

	localLoc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal(err)
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

	for _, t := range times {
		if ctx.String("from") != "" && from.After(t.Created) {
			continue
		}
		if ctx.String("until") != "" && until.Before(t.Created) {
			continue
		}

		totalDuration += t.Time

		outputValues = append(
			outputValues,
			[]string{
				strconv.FormatInt(t.ID, 10),
				t.Created.In(localLoc).Format("2006-01-02 15:04:05"),
				"#" + strconv.FormatInt(t.IssueID, 10),
				strconv.FormatInt(t.UserID, 10),
				time.Duration(1e9 * t.Time).String(),
			},
		)
	}

	if ctx.Bool("total") {
		outputValues = append(outputValues, []string{
			"TOTAL", "", "", "", time.Duration(1e9 * totalDuration).String(),
		})
	}

	Output(outputValue, headers, outputValues)
	return nil
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

	if len(ctx.Args()) < 2 {
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
