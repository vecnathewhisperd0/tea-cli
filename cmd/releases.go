// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReleases represents to login a gitea server.
// ToDo: ReleaseDetails
var CmdReleases = cli.Command{
	Name:        "release",
	Aliases:     []string{"releases"},
	Usage:       "Manage releases",
	Description: "Manage releases",
	Action:      runReleasesList,
	Subcommands: []*cli.Command{
		&CmdReleaseList,
		&CmdReleaseCreate,
		&CmdReleaseDelete,
		&CmdReleaseEdit,
	},
	Flags: AllDefaultFlags,
}

func getReleaseByTag(owner, repo, tag string, client *gitea.Client) (*gitea.Release, error) {
	rl, _, err := client.ListReleases(owner, repo, gitea.ListReleasesOptions{})
	if err != nil {
		return nil, err
	}
	if len(rl) == 0 {
		fmt.Println("Repo does not have any release")
		return nil, nil
	}
	for _, r := range rl {
		if r.TagName == tag {
			return r, nil
		}
	}
	fmt.Println("Release tag does not exist")
	return nil, nil
}
