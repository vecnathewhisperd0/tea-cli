// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/task"

	"code.gitea.io/sdk/gitea"
	"github.com/AlecAivazis/survey/v2"
)

// CreateIssue interactively creates an issue
func CreateIssue(login *config.Login, owner, repo string) error {
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	var opts gitea.CreateIssueOption

	// TODO: pass client instance down for speed
	templates, _, err := login.Client().GetIssueTemplates(owner, repo)

	if err == nil && len(templates) != 0 {
		// FIXME perf don't do this twice.
		selectableChan := make(chan (issueSelectables), 1)
		go fetchIssueSelectables(login, owner, repo, selectableChan)

		var (
			templateNames   = make([]string, len(templates))
			templatesByName = map[string]*gitea.IssueTemplate{}
		)
		for i, t := range templates {
			name := t.Name + "\t " + t.About
			templatesByName[name] = t
			templateNames[i] = name
		}

		if selectedTemplate, err := promptSelect("Use issue template?", templateNames, "", "[none]"); err != nil {
			return err
		} else if selectedTemplate != "[none]" {
			// wait until selectables are fetched
			selectables := <-selectableChan
			if selectables.Err != nil {
				return selectables.Err
			}
			opts, err = promptForIssueTemplate(templatesByName[selectedTemplate], selectables.LabelMap)
			if err != nil {
				return err
			}
		}
	}

	if err := promptIssueProperties(login, owner, repo, &opts); err != nil {
		return err
	}

	return task.CreateIssue(login, owner, repo, opts)
}

func promptForIssueTemplate(t *gitea.IssueTemplate, labels map[string]int64) (gitea.CreateIssueOption, error) {
	// map label names to IDs
	var labelIDs = make([]int64, len(t.IssueLabels))
	for i, l := range t.IssueLabels {
		labelIDs[i] = labels[l]
	}

	opts := gitea.CreateIssueOption{
		Title:  t.IssueTitle,
		Ref:    t.IssueRef,
		Labels: labelIDs,
	}

	if t.IsForm() {
		var responses []string
		for _, f := range t.Form {
			var (
				promptOpts survey.AskOpt
				prompt     survey.Prompt
				input      string
				// TODO: reuse code from https://github.com/go-gitea/gitea/blob/main/modules/issue/template/template.go#L256-L306
				responseFormatter = func(in string) string { return in }
			)

			if f.Validations.Required {
				promptOpts = survey.WithValidator(survey.Required)
			}

			switch f.Type {
			case gitea.IssueFormElementMarkdown:
				print.OutputMarkdown(f.Attributes.Value, "")
				continue // not emitted to issue body
			case gitea.IssueFormElementInput:
				prompt = &survey.Input{
					Message: f.Attributes.Label,
					Help:    f.Attributes.Description,
					Default: f.Attributes.Value,
				}
				responseFormatter = func(answer string) string {
					return fmt.Sprintf("#### %s\n%s", f.Attributes.Label, answer)
				}
			}

			if prompt != nil {
				if err := survey.AskOne(prompt, &input, promptOpts); err != nil {
					return opts, err
				}
				responses = append(responses, responseFormatter(input))
			}
		}
	} else {
		// TODO: don't prompt for description again, if we did here.
		prompt := NewMultiline(Multiline{
			Message:   "Issue Description",
			Default:   t.MarkdownContent,
			Syntax:    "md",
			UseEditor: config.GetPreferences().Editor,
		})
		if err := survey.AskOne(prompt, &opts.Body); err != nil {
			return opts, err
		}
	}

	return opts, nil
}

func promptIssueProperties(login *config.Login, owner, repo string, o *gitea.CreateIssueOption) error {
	var milestoneName string
	var labels []string
	var err error

	selectableChan := make(chan (issueSelectables), 1)
	go fetchIssueSelectables(login, owner, repo, selectableChan)

	// title
	promptOpts := survey.WithValidator(survey.Required)
	promptI := &survey.Input{Message: "Issue title:", Default: o.Title}
	if err = survey.AskOne(promptI, &o.Title, promptOpts); err != nil {
		return err
	}

	// description
	promptD := NewMultiline(Multiline{
		Message:   "Issue description:",
		Default:   o.Body,
		Syntax:    "md",
		UseEditor: config.GetPreferences().Editor,
	})
	if err = survey.AskOne(promptD, &o.Body); err != nil {
		return err
	}

	// wait until selectables are fetched
	selectables := <-selectableChan
	if selectables.Err != nil {
		return selectables.Err
	}

	// skip remaining props if we don't have permission to set them
	if !selectables.Repo.Permissions.Push {
		return nil
	}

	// assignees
	if o.Assignees, err = promptMultiSelect("Assignees:", selectables.Assignees, "[other]"); err != nil {
		return err
	}

	// milestone
	if len(selectables.MilestoneList) != 0 {
		if milestoneName, err = promptSelect("Milestone:", selectables.MilestoneList, selectables.MilestoneMapInv[o.Milestone], "[none]"); err != nil {
			return err
		}
		o.Milestone = selectables.MilestoneMap[milestoneName]
	}

	// labels
	if len(selectables.LabelList) != 0 {
		for _, l := range o.Labels {
			labels = append(labels, selectables.LabelMapInv[l])
		}
		promptL := &survey.MultiSelect{Message: "Labels:", Options: selectables.LabelList, VimMode: true, Default: labels}
		if err := survey.AskOne(promptL, &labels); err != nil {
			return err
		}
		o.Labels = make([]int64, len(labels))
		for i, l := range labels {
			o.Labels[i] = selectables.LabelMap[l]
		}
	}

	// deadline
	if o.Deadline, err = promptDatetime("Due date:"); err != nil {
		return err
	}

	return nil
}

type issueSelectables struct {
	Repo            *gitea.Repository
	Assignees       []string
	MilestoneList   []string
	MilestoneMap    map[string]int64
	MilestoneMapInv map[int64]string
	LabelList       []string
	LabelMap        map[string]int64
	LabelMapInv     map[int64]string
	Err             error
}

func fetchIssueSelectables(login *config.Login, owner, repo string, done chan issueSelectables) {
	// TODO PERF make these calls concurrent
	r := issueSelectables{}
	c := login.Client()

	r.Repo, _, r.Err = c.GetRepo(owner, repo)
	if r.Err != nil {
		done <- r
		return
	}
	// we can set the following properties only if we have write access to the repo
	// so we fastpath this if not.
	if !r.Repo.Permissions.Push {
		done <- r
		return
	}

	assignees, _, err := c.GetAssignees(owner, repo)
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.Assignees = make([]string, len(assignees))
	for i, u := range assignees {
		r.Assignees[i] = u.UserName
	}

	milestones, _, err := c.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.MilestoneMap = make(map[string]int64)
	r.MilestoneMapInv = make(map[int64]string)
	r.MilestoneList = make([]string, len(milestones))
	for i, m := range milestones {
		r.MilestoneMap[m.Title] = m.ID
		r.MilestoneMapInv[m.ID] = m.Title
		r.MilestoneList[i] = m.Title
	}

	labels, _, err := c.ListRepoLabels(owner, repo, gitea.ListLabelsOptions{
		ListOptions: gitea.ListOptions{Page: -1},
	})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.LabelMap = make(map[string]int64)
	r.LabelMapInv = make(map[int64]string)
	r.LabelList = make([]string, len(labels))
	for i, l := range labels {
		r.LabelMap[l.Name] = l.ID
		r.LabelMapInv[l.ID] = l.Name
		r.LabelList[i] = l.Name
	}

	done <- r
}
