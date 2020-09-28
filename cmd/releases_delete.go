// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/tea/modules/intern"

	"github.com/urfave/cli/v2"
)

// CmdReleaseDelete represents a sub command of Release to delete a release
var CmdReleaseDelete = cli.Command{
	Name:        "delete",
	Usage:       "Delete a release",
	Description: `Delete a release`,
	ArgsUsage:   "<release tag>",
	Action:      runReleaseDelete,
	Flags:       AllDefaultFlags,
}

func runReleaseDelete(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)
	client := login.Client()

	tag := ctx.Args().First()
	if len(tag) == 0 {
		fmt.Println("Release tag needed to delete")
		return nil
	}

	release, err := getReleaseByTag(owner, repo, tag, client)
	if err != nil {
		return err
	}
	if release == nil {
		return nil
	}

	_, err = client.DeleteRelease(owner, repo, release.ID)
	return err
}
