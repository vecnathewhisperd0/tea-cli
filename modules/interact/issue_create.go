// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
	"github.com/araddon/dateparse"
)

// CreateIssue interactively creates an issue
func CreateIssue(login *config.Login, owner, repo string) error {
	var title, description, dueDate, assignees, milestone, labels string
	var deadline *time.Time
	var msID int64
	var labelIDs []int64

	// owner, repo
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	// title
	promptOpts := survey.WithValidator(survey.Required)
	promptI := &survey.Input{Message: "Issue title:"}
	if err := survey.AskOne(promptI, &title, promptOpts); err != nil {
		return err
	}

	// description
	promptD := &survey.Multiline{Message: "Issue description:"}
	if err := survey.AskOne(promptD, &description); err != nil {
		return err
	}

	// assignees // TODO: add suggestions
	promptA := &survey.Input{Message: "Assignees:"}
	if err := survey.AskOne(promptA, &assignees); err != nil {
		return err
	}

	// milestone // TODO: add suggestions
	promptM := &survey.Input{Message: "Milestone:"}
	if err := survey.AskOne(promptM, &milestone); err != nil {
		return err
	}

	// labels // TODO: add suggestions
	promptL := &survey.Input{Message: "Labels:"}
	if err := survey.AskOne(promptL, &labels); err != nil {
		return err
	}

	// deadline
	promptI = &survey.Input{Message: "Due date [no due date]:"}
	err = survey.AskOne(
		promptI,
		&dueDate,
		survey.WithValidator(func(input interface{}) error {
			if str, ok := input.(string); ok {
				if len(str) == 0 {
					return nil
				}
				t, err := dateparse.ParseAny(str)
				if err != nil {
					return err
				}
				deadline = &t
			} else {
				return fmt.Errorf("invalid result type")
			}
			return nil
		}),
	)

	// resolve IDs
	client := login.Client()
	if len(milestone) != 0 {
		ms, _, err := client.GetMilestoneByName(owner, repo, milestone)
		if err != nil {
			return fmt.Errorf("Milestone '%s' not found", milestone)
		}
		msID = ms.ID
	}
	if len(labels) != 0 {
		labelIDs, err = task.ResolveLabelNames(client, owner, repo, strings.Split(labels, ","))
		if err != nil {
			return err
		}
	}

	return task.CreateIssue(
		login,
		owner,
		repo,
		gitea.CreateIssueOption{
			Title:     title,
			Body:      description,
			Deadline:  deadline,
			Assignees: strings.Split(assignees, ","),
			Milestone: msID,
			Labels:    labelIDs,
		},
	)
}
