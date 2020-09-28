// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/tea/modules/intern"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReleaseList represents a sub command of Release to list releases
var CmdReleaseList = cli.Command{
	Name:        "ls",
	Usage:       "List Releases",
	Description: "List Releases",
	Action:      runReleasesList,
	Flags: append([]cli.Flag{
		&PaginationPageFlag,
		&PaginationLimitFlag,
	}, AllDefaultFlags...),
}

func runReleasesList(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	releases, _, err := login.Client().ListReleases(owner, repo, gitea.ListReleasesOptions{ListOptions: getListOptions(ctx)})
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Tag-Name",
		"Title",
		"Published At",
		"Status",
		"Tar URL",
	}

	var values [][]string

	if len(releases) == 0 {
		Output(globalOutputValue, headers, values)
		return nil
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
				release.PublishedAt.Format("2006-01-02 15:04:05"),
				status,
				release.TarURL,
			},
		)
	}
	Output(globalOutputValue, headers, values)

	return nil
}
