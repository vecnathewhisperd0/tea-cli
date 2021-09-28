// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// MilestoneDetails print an milestone formatted to stdout
func MilestoneDetails(milestone *gitea.Milestone) {
	fmt.Printf("%s\n",
		milestone.Title,
	)
	if len(milestone.Description) != 0 {
		fmt.Printf("\n%s\n", milestone.Description)
	}
	if milestone.Deadline != nil && !milestone.Deadline.IsZero() {
		fmt.Printf("\nDeadline: %s\n", FormatTime(*milestone.Deadline))
	}
}

// MilestonesList prints a listing of milestones
func MilestonesList(news []*gitea.Milestone, output string, fields []string) {
	var printables = make([]printable, len(news))
	for i, x := range news {
		printables[i] = &printableMilestone{x}
	}
	t := tableFromItems(fields, printables)
	t.sort(0, true)
	t.print(output)
}

// MilestoneFields are all available fields to print with MilestonesList
var MilestoneFields = []string{
	"title",
	"state",
	"open items",
	"closed items",
	"open/closed issues",
	"due date",
	"description",
	"created",
	"updated",
	"closed",
	"id",
}

type printableMilestone struct {
	*gitea.Milestone
}

func (m printableMilestone) FormatField(field string) string {
	switch field {
	case "title":
		return m.Title
	case "state":
		return string(m.State)
	case "open items":
		return fmt.Sprintf("%d", m.OpenIssues)
	case "closed items":
		return fmt.Sprintf("%d", m.ClosedIssues)
	case "open/closed issues": // for backwards compatibility
		return fmt.Sprintf("%d/%d", m.OpenIssues, m.ClosedIssues)
	case "deadline", "due date":
		if m.Deadline != nil && !m.Deadline.IsZero() {
			return FormatTime(*m.Deadline)
		}
	case "id":
		return fmt.Sprintf("%d", m.ID)
	case "description":
		return m.Description
	case "created":
		return FormatTime(m.Created)
	case "updated":
		if m.Updated != nil {
			return FormatTime(*m.Updated)
		}
	case "closed":
		if m.Closed != nil {
			return FormatTime(*m.Closed)
		}
	}
	return ""
}
