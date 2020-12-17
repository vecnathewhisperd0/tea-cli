// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"time"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/AlecAivazis/survey/v2"
)

// CreateMilestone interactively creates a milestone
func CreateMilestone(login *config.Login, owner, repo string, deadline *time.Time, state gitea.StateType) error {
	var title, description, dueDate string

	// owner, repo
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	// title
	promptOpts := survey.WithValidator(survey.Required)
	promptI := &survey.Input{Message: "Milestone title:"}
	if err := survey.AskOne(promptI, &title, promptOpts); err != nil {
		return err
	}

	// description
	promptM := &survey.Multiline{Message: "Milestone description:"}
	if err := survey.AskOne(promptM, &description); err != nil {
		return err
	}

	// deadline
	promptI = &survey.Input{Message: "Milestone deadline [no due date]:"}
	if err := survey.AskOne(promptI, &dueDate, nil); err != nil {
		return err
	}
	if dueDate != "" {
		deadline, err = utils.GetIso8601Date(dueDate)
		if err != nil {
			return err
		}
	} else {
		deadline = nil
	}

	return task.CreateMilestone(
		login,
		owner,
		repo,
		title,
		description,
		deadline,
		state)
}
