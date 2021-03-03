// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"time"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
)

// CreateIssue interactively creates an issue
func CreateIssue(login *config.Login, owner, repo string) error {
	var title, description, milestone string
	var assignees, labels []string
	var deadline *time.Time

	// owner, repo
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	selectableChan := make(chan (issueSelectables), 1)
	go fetchIssueSelectables(login, owner, repo, selectableChan)

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

	// wait until selectables are fetched
	selectables := <-selectableChan
	if selectables.Err != nil {
		return selectables.Err
	}

	// assignees
	if assignees, err = promptMultiSelect("Assignees:", selectables.Collaborators, "[other]"); err != nil {
		return err
	}

	// milestone
	if milestone, err = promptSelect("Milestone:", selectables.MilestoneList, "", "[none]"); err != nil {
		return err
	}

	// labels
	promptL := &survey.MultiSelect{Message: "Labels:", Options: selectables.LabelList, VimMode: true}
	if err := survey.AskOne(promptL, &labels); err != nil {
		return err
	}
	labelIDs := make([]int64, len(labels))
	for i, l := range labels {
		labelIDs[i] = selectables.LabelMap[l]
	}

	// deadline
	if deadline, err = promptDatetime("Due date:"); err != nil {
		return err
	}

	return task.CreateIssue(
		login,
		owner,
		repo,
		gitea.CreateIssueOption{
			Title:     title,
			Body:      description,
			Deadline:  deadline,
			Assignees: assignees,
			Milestone: selectables.MilestoneMap[milestone],
			Labels:    labelIDs,
		},
	)
}

type issueSelectables struct {
	Collaborators []string
	MilestoneList []string
	MilestoneMap  map[string]int64
	LabelList     []string
	LabelMap      map[string]int64
	Err           error
}

func fetchIssueSelectables(login *config.Login, owner, repo string, done chan issueSelectables) {
	// TODO PERF make these calls concurrent
	r := issueSelectables{}
	c := login.Client()

	// FIXME: this should ideally be ListAssignees(), https://github.com/go-gitea/gitea/issues/14856
	colabs, _, err := c.ListCollaborators(owner, repo, gitea.ListCollaboratorsOptions{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.Collaborators = make([]string, len(colabs)+1)
	r.Collaborators[0] = login.User
	for i, u := range colabs {
		r.Collaborators[i+1] = u.UserName
	}

	milestones, _, err := c.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.MilestoneMap = make(map[string]int64)
	r.MilestoneList = make([]string, len(milestones))
	for i, m := range milestones {
		r.MilestoneMap[m.Title] = m.ID
		r.MilestoneList[i] = m.Title
	}

	labels, _, err := c.ListRepoLabels(owner, repo, gitea.ListLabelsOptions{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.LabelMap = make(map[string]int64)
	r.LabelList = make([]string, len(labels))
	for i, l := range labels {
		r.LabelMap[l.Name] = l.ID
		r.LabelList[i] = l.Name
	}

	done <- r
}
