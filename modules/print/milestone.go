// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

func MilestoneDetails(milestone *gitea.Milestone) {
	fmt.Printf("%s\n",
		milestone.Title,
	)
	if len(milestone.Description) != 0 {
		fmt.Printf("\n%s\n", milestone.Description)
	}
	if milestone.Deadline != nil && !milestone.Deadline.IsZero() {
		fmt.Printf("\nDeadline: %s\n", milestone.Deadline.Format("2006-01-02 15:04:05"))
	}
}
