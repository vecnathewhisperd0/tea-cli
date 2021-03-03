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
	"code.gitea.io/tea/modules/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/araddon/dateparse"
)

const nilVal = "[none]"
const customVal = "[other]"

// CreateIssue interactively creates an issue
func CreateIssue(login *config.Login, owner, repo string) error {
	var title, description, dueDate, milestone string
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
	promptA := &survey.MultiSelect{Message: "Assignees:", Options: selectables.Collaborators, VimMode: true}
	if err := survey.AskOne(promptA, &assignees); err != nil {
		return err
	}
	// check for custom value & prompt again with text input
	// HACK until https://github.com/AlecAivazis/survey/issues/339 is implemented
	if otherIndex := utils.IndexOf(assignees, customVal); otherIndex != -1 {
		var customAssignees string
		promptA := &survey.Input{Message: "Assignees:", Help: "comma separated usernames"}
		if err := survey.AskOne(promptA, &customAssignees); err != nil {
			return err
		}
		assignees = append(assignees[:otherIndex], assignees[otherIndex+1:]...)
		assignees = append(assignees, strings.Split(customAssignees, ",")...)
	}

	// milestone
	promptM := &survey.Select{Message: "Milestone:", Options: selectables.MilestoneList, VimMode: true, Default: nilVal}
	if err := survey.AskOne(promptM, &milestone); err != nil {
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
	r.Collaborators = make([]string, len(colabs)+2)
	r.Collaborators[0] = login.User
	r.Collaborators[1] = customVal
	for i, u := range colabs {
		r.Collaborators[i+2] = u.UserName
	}

	milestones, _, err := c.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.MilestoneMap = make(map[string]int64)
	r.MilestoneList = make([]string, len(milestones)+1)
	r.MilestoneList[0] = nilVal
	r.MilestoneMap[nilVal] = 0
	for i, m := range milestones {
		r.MilestoneMap[m.Title] = m.ID
		r.MilestoneList[i+1] = m.Title
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
