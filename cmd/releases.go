// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli/v2"
)

// CmdReleases represents to login a gitea server.
var CmdReleases = cli.Command{
	Name:        "releases",
	Usage:       "Create releases",
	Description: `Create releases`,
	Action:      runReleases,
	Subcommands: []*cli.Command{
		&CmdReleaseCreate,
		&CmdReleaseDelete,
	},
	Flags: AllDefaultFlags,
}

func runReleases(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	releases, _, err := login.Client().ListReleases(owner, repo, gitea.ListReleasesOptions{})
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
		Output(outputValue, headers, values)
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
	Output(outputValue, headers, values)

	return nil
}

// CmdReleaseCreate represents a sub command of Release to create release
var CmdReleaseCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a release",
	Description: `Create a release`,
	Action:      runReleaseCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "tag",
			Usage: "Tag name",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: "Target refs, branch name or commit id",
		},
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "Release title",
		},
		&cli.StringFlag{
			Name:    "note",
			Aliases: []string{"n"},
			Usage:   "Release notes",
		},
		&cli.BoolFlag{
			Name:    "draft",
			Aliases: []string{"d"},
			Usage:   "Is a draft",
		},
		&cli.BoolFlag{
			Name:    "prerelease",
			Aliases: []string{"p"},
			Usage:   "Is a pre-release",
		},
		&cli.StringSliceFlag{
			Name:    "asset",
			Aliases: []string{"a"},
			Usage:   "List of files to attach",
		},
	}, LoginRepoFlags...),
}

func runReleaseCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	release, resp, err := login.Client().CreateRelease(owner, repo, gitea.CreateReleaseOption{
		TagName:      ctx.String("tag"),
		Target:       ctx.String("target"),
		Title:        ctx.String("title"),
		Note:         ctx.String("note"),
		IsDraft:      ctx.Bool("draft"),
		IsPrerelease: ctx.Bool("prerelease"),
	})

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusConflict {
			fmt.Println("error: There already is a release for this tag")
			return nil
		}
		log.Fatal(err)
	}

	for _, asset := range ctx.StringSlice("asset") {
		var file *os.File

		if file, err = os.Open(asset); err != nil {
			log.Fatal(err)
		}

		filePath := filepath.Base(asset)

		if _, _, err = login.Client().CreateReleaseAttachment(owner, repo, release.ID, file, filePath); err != nil {
			file.Close()
			log.Fatal(err)
		}

		file.Close()
	}

	return nil
}

// CmdReleaseDelete represents a sub command of Release to delete a release
var CmdReleaseDelete = cli.Command{
	Name:        "delete",
	Usage:       "Delete a release",
	Description: `Delete a release`,
	ArgsUsage:   "[<release tag>]",
	Action:      runReleaseDelete,
	Flags:       LoginRepoFlags,
}

func runReleaseDelete(ctx *cli.Context) error {
	login, owner, repo := initCommand()
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
