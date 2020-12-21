// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"strconv"

	"code.gitea.io/sdk/gitea"
)

// TrackedTimesList print list of tracked times to stdout
func TrackedTimesList(times []*gitea.TrackedTime, outputType string, printTotal bool) {
	tab := tableWithHeader(
		"Created",
		"Issue",
		"User",
		"Duration",
	)
	var totalDuration int64

	for _, t := range times {
		totalDuration += t.Time
		tab.addRow(
			FormatTime(t.Created),
			"#"+strconv.FormatInt(t.Issue.Index, 10),
			t.UserName,
			formatDuration(t.Time, outputType),
		)
	}

	if printTotal {
		tab.addRow("TOTAL", "", "", formatDuration(totalDuration, outputType))
	}
	tab.print(outputType)
}
