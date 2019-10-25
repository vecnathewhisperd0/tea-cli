// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"os"
	"strconv"
	"time"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdLabels represents to operate repositories' labels.
var CmdTrackedTimes = cli.Command{
	Name:        "times",
	Usage:       `Operate on tracked times of the repository's issues`,
	Description: `Operate on tracked times of the repository's issues`,
	Action:      runTrackedTimes,
	Subcommands: []cli.Command{
		// CmdTimeAdd,
	},
	Flags: AllDefaultFlags,
}

func runTrackedTimes(ctx *cli.Context) error {
	login, owner, repo := initCommand()

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

	if len(os.Args) == 3 {
		times, err = login.Client().GetUserTrackedTimes(owner, repo, os.Args[2])
	} else {
		times, err = login.Client().GetRepoTrackedTimes(owner, repo)
	}

	if err != nil {
		log.Fatal(err)
	}

	localLoc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range times {

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
	Output(outputValue, headers, outputValues)

	return nil
}
