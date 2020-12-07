// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"code.gitea.io/tea/cmd/flags"

	"code.gitea.io/sdk/gitea"
)

// ReleasesList prints a listing of releases
func ReleasesList(releases []*gitea.Release) {
	var values [][]string
	headers := []string{
		"Tag-Name",
		"Title",
		"Published At",
		"Status",
		"Tar URL",
	}

	if len(releases) == 0 {
		OutputList(flags.GlobalOutputValue, headers, values)
		return
	}

	for _, release := range releases {
		status := "released"
		if release.IsDraft {
			status = "draft"
		} else if release.IsPrerelease {
			status = "prerelease"
		}
		values = append(
			values,
			[]string{
				release.TagName,
				release.Title,
				FormatTime(release.PublishedAt),
				status,
				release.TarURL,
			},
		)
	}

	OutputList(flags.GlobalOutputValue, headers, values)
}
