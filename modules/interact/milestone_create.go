// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"

	"code.gitea.io/sdk/gitea"
	"github.com/AlecAivazis/survey/v2"
)

// CreateMilestone interactively creates a milestone
func CreateMilestone(login *config.Login, owner, repo string, state gitea.StateType) error {
	var title, description string

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

	return task.CreateMilestone(
		login,
		owner,
		repo,
		title,
		description,
		state)
}
