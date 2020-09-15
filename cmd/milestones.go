// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestones represents to operate repositories milestones.
var CmdMilestones = cli.Command{
	Name:        "milestones",
	Aliases:     []string{"ms"},
	Usage:       "List and create milestones",
	Description: `List and create milestones`,
	ArgsUsage:   "[<milestone name/id>]",
	Action:      runMilestones,
	Subcommands: []*cli.Command{
		&CmdMilestonesList,
		&CmdMilestonesCreate,
		&CmdMilestonesClose,
		&CmdMilestonesRemove,
		&CmdMilestonesReopen,
	},
	Flags: AllDefaultFlags,
}

// CmdMilestonesList represents a sub command of milestones to list milestones
var CmdMilestonesList = cli.Command{
	Name:        "ls",
	Usage:       "List milestones of the repository",
	Description: `List milestones of the repository`,
	Action:      runMilestonesList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by milestone state (all|open|closed)",
			DefaultText: "open",
		},
	}, AllDefaultFlags...),
}

func runMilestones(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runMilestoneDetail(ctx, ctx.Args().First())
	}
	return runMilestonesList(ctx)
}

func runMilestoneDetail(ctx *cli.Context, value string) error {
	login, owner, repo := initCommand()
	client := login.Client()

	mileID, err := getMilestoneID(client, owner, repo, value)
	if err != nil {
		return err
	}
	milestone, err := client.GetMilestone(owner, repo, mileID)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n",
		milestone.Title,
	)
	if len(milestone.Description) != 0 {
		fmt.Printf("\n%s\n", milestone.Description)
	}
	if milestone.Deadline != nil && !milestone.Deadline.IsZero() {
		fmt.Printf("\nDeadline: %s\n", milestone.Deadline.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func runMilestonesList(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "closed":
		state = gitea.StateClosed
	}

	milestones, err := login.Client().ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{
		State: state,
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"ID",
		"Title",
	}
	if state == gitea.StateAll {
		headers = append(headers, "State")
	}
	headers = append(headers,
		"Open/Closed Issues",
		"DueDate",
	)

	var values [][]string

	for _, m := range milestones {
		var deadline = ""

		if m.Deadline != nil && !m.Deadline.IsZero() {
			deadline = m.Deadline.Format("2006-01-02 15:04:05")
		}

		item := []string{
			fmt.Sprintf("%d", m.ID),
			m.Title,
		}
		if state == gitea.StateAll {
			item = append(item, string(m.State))
		}
		item = append(item,
			fmt.Sprintf("%d/%d", m.OpenIssues, m.ClosedIssues),
			deadline,
		)

		values = append(values, item)
	}
	Output(outputValue, headers, values)

	return nil
}

// CmdMilestonesCreate represents a sub command of milestones to create milestone
var CmdMilestonesCreate = cli.Command{
	Name:        "create",
	Usage:       "Create an milestone on repository",
	Description: `Create an milestone on repository`,
	Action:      runMilestonesCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "milestone title to create",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "milestone description to create",
		},
		&cli.StringFlag{
			Name:        "state",
			Usage:       "set milestone state (default is open)",
			DefaultText: "open",
		},
	}, AllDefaultFlags...),
}

func runMilestonesCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	title := ctx.String("title")
	if len(title) == 0 {
		fmt.Printf("Title is required\n")
		return nil
	}

	state := gitea.StateOpen
	if ctx.String("state") == "closed" {
		state = gitea.StateClosed
	}

	mile, err := login.Client().CreateMilestone(owner, repo, gitea.CreateMilestoneOption{
		Title:       title,
		Description: ctx.String("description"),
		State:       state,
	})
	if err != nil {
		log.Fatal(err)
	}

	return runMilestoneDetail(ctx, fmt.Sprintf("%d", mile.ID))
}

// CmdMilestonesClose represents a sub command of milestones to close an milestone
var CmdMilestonesClose = cli.Command{
	Name:        "close",
	Usage:       "Change state of an milestone to 'closed'",
	Description: `Change state of an milestone to 'closed'`,
	ArgsUsage:   "<milestone name/id>",
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("force") {
			return deleteMilestone(ctx)
		}
		return editMilestoneStatus(ctx, true)
	},
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "delete milestone",
		},
	}, AllDefaultFlags...),
}

func editMilestoneStatus(ctx *cli.Context, close bool) error {
	login, owner, repo := initCommand()
	client := login.Client()

	mileID, err := getMilestoneID(client, owner, repo, ctx.Args().First())
	if err != nil {
		return err
	}
	if mileID != 0 {

		state := string(gitea.StateOpen)
		if close {
			state = string(gitea.StateClosed)
		}
		_, err = client.EditMilestone(owner, repo, mileID, gitea.EditMilestoneOption{
			State: &state,
		})
	}
	return err
}

// CmdMilestonesRemove represents a sub command of milestones to delete an milestone
var CmdMilestonesRemove = cli.Command{
	Name:        "remove",
	Aliases:     []string{"rm"},
	Usage:       "delete a milestone",
	Description: "delete a milestone",
	ArgsUsage:   "<milestone name/id>",
	Action:      deleteMilestone,
	Flags:       AllDefaultFlags,
}

func deleteMilestone(ctx *cli.Context) error {
	login, owner, repo := initCommand()
	client := login.Client()

	mileID, err := getMilestoneID(client, owner, repo, ctx.Args().First())
	if err != nil {
		return err
	}

	return client.DeleteMilestone(owner, repo, mileID)
}

// CmdMilestonesReopen represents a sub command of milestones to open an milestone
var CmdMilestonesReopen = cli.Command{
	Name:        "reopen",
	Aliases:     []string{"open"},
	Usage:       "Change state of an milestone to 'open'",
	Description: `Change state of an milestone to 'open'`,
	ArgsUsage:   "<milestone name/id>",
	Action: func(ctx *cli.Context) error {
		return editMilestoneStatus(ctx, false)
	},
	Flags: AllDefaultFlags,
}

// getMilestoneID TODO: delete it and use sdk feature when v0.13.0 is released
func getMilestoneID(client *gitea.Client, owner, repo, nameOrID string) (int64, error) {
	if match, err := regexp.MatchString("^\\d+$", nameOrID); err == nil && match {
		return strconv.ParseInt(nameOrID, 10, 64)
	}
	i := 0
	for {
		i++
		miles, err := client.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{
			ListOptions: gitea.ListOptions{
				Page: i,
			},
			State: "all",
		})
		if err != nil {
			return 0, err
		}
		if len(miles) == 0 {
			return 0, nil
		}
		for _, m := range miles {
			if strings.ToLower(strings.TrimSpace(m.Title)) == strings.ToLower(strings.TrimSpace(nameOrID)) {
				return m.ID, nil
			}
		}
	}
}
