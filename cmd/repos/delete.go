// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repos

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v2"

	"github.com/AlecAivazis/survey/v2"
)

// CmdRepoFork represents a sub command of repos to fork an existing repo
var CmdRepoRm = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "Delete an existing repository",
	Description: "Removes a repository from Create a repository from an existing repo",
	ArgsUsage:   " ", // command does not accept arguments
	Action:      runRepoDelete,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{""},
			Required: true,
			Usage:    "name of the repo",
		},
		&cli.StringFlag{
			Name:     "owner",
			Aliases:  []string{"O"},
			Required: false,
			Usage:    "owner of the repo",
		},
		&cli.BoolFlag{
			Name:     "force",
			Aliases:  []string{"f"},
			Required: false,
			Value:    false,
			Usage:    "Force the deletion and don't ask for confirmation",
		},
	}, flags.LoginOutputFlags...),
}

func runRepoDelete(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)

	client := ctx.Login.Client()

	var owner string
	if ctx.IsSet("owner") {
		owner = ctx.String("owner")

	} else {
		owner = ctx.Login.User
	}

	repo_name := ctx.String("name")

	repo_slug := fmt.Sprintf("%s/%s", owner, repo_name)

	if !ctx.Bool("force") {
		var entered_repo_slug string
		promptRepoName := &survey.Input{
			Message: fmt.Sprintf("Confirm the deletion of the repository '%s' by typing its name: ", repo_slug),
		}
		if err := survey.AskOne(promptRepoName, &entered_repo_slug, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		if entered_repo_slug != repo_slug {
			return fmt.Errorf("Entered wrong repository name '%s', expected '%s'", entered_repo_slug, repo_slug)
		}
	}

	_, err := client.DeleteRepo(owner, repo_name)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully deleted %s/%s\n", owner, repo_name)
	return nil
}
