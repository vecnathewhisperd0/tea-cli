// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "show repositories details",
	Description: `Operate with repositories`,
	ArgsUsage:   "[<repo owner>/<repo name>]",
	Action:      runRepos,
	Subcommands: []*cli.Command{
		&CmdReposList,
		&CmdRepoCreate,
	},
	Flags: LoginOutputFlags,
}

// CmdReposList represents a sub command of repos to list them
var CmdReposList = cli.Command{
	Name:        "ls",
	Usage:       "List available repositories",
	Description: `List available repositories`,
	Action:      runReposList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "mode",
			Aliases:  []string{"m"},
			Required: false,
			Usage:    "Filter by mode: fork, mirror, source",
		},
		&cli.StringFlag{
			Name:     "org",
			Required: false,
			Usage:    "Filter by organization",
		},
		&cli.StringFlag{
			Name:     "user",
			Aliases:  []string{"u"},
			Required: false,
			Usage:    "Filter by user",
		},
	}, LoginOutputFlags...),
}

// CmdRepoCreate represents a sub command of repos to create one
var CmdRepoCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create a repository",
	Description: "Create a repository",
	Action:      runRepoCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{""},
			Required: true,
			Usage:    "name of new repo",
		},
		&cli.StringFlag{
			Name:     "owner",
			Aliases:  []string{"O"},
			Required: false,
			Usage:    "name of repo owner",
		},
		&cli.BoolFlag{
			Name:     "private",
			Required: false,
			Value:    false,
			Usage:    "make repo private",
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"desc"},
			Required: false,
			Usage:    "add description to repo",
		},
		&cli.BoolFlag{
			Name:     "init",
			Required: false,
			Value:    false,
			Usage:    "initialize repo",
		},
		&cli.StringFlag{
			Name:     "labels",
			Required: false,
			Usage:    "name of label set to add",
		},
		&cli.StringFlag{
			Name:     "gitignores",
			Aliases:  []string{"git"},
			Required: false,
			Usage:    "list of gitignore templates (need --init)",
		},
		&cli.StringFlag{
			Name:     "license",
			Required: false,
			Usage:    "add license (need --init)",
		},
		&cli.StringFlag{
			Name:     "readme",
			Required: false,
			Usage:    "use readme template (need --init)",
		},
		&cli.StringFlag{
			Name:     "branch",
			Required: false,
			Usage:    "use custom default branch (need --init)",
		},
	}, LoginOutputFlags...),
}

func runRepos(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runRepoDetail(ctx, ctx.Args().First())
	}
	return runReposList(ctx)
}

// runReposList list repositories
func runReposList(ctx *cli.Context) error {
	login := initCommandLoginOnly()

	mode := ctx.String("mode")
	org := ctx.String("org")
	user := ctx.String("user")

	var rps []*gitea.Repository
	var err error

	// TODO: on sdk v0.13.0 release, switch to SearchRepos()
	// Note: user filter can be used as org filter too
	if org != "" {
		rps, err = login.Client().ListOrgRepos(org, gitea.ListOrgReposOptions{})
	} else if user != "" {
		rps, err = login.Client().ListUserRepos(user, gitea.ListReposOptions{})
	} else {
		rps, err = login.Client().ListMyRepos(gitea.ListReposOptions{})
	}
	if err != nil {
		log.Fatal(err)
	}

	var repos []*gitea.Repository
	if mode == "" {
		repos = rps
	} else if mode == "fork" {
		for _, rp := range rps {
			if rp.Fork == true {
				repos = append(repos, rp)
			}
		}
	} else if mode == "mirror" {
		for _, rp := range rps {
			if rp.Mirror == true {
				repos = append(repos, rp)
			}
		}
	} else if mode == "source" {
		for _, rp := range rps {
			if rp.Mirror != true && rp.Fork != true {
				repos = append(repos, rp)
			}
		}
	} else {
		log.Fatal("Unknown mode: ", mode, "\nUse one of the following:\n- fork\n- mirror\n- source\n")
		return nil
	}

	if len(rps) == 0 {
		log.Fatal("No repositories found", rps)
		return nil
	}

	headers := []string{
		"Name",
		"Type",
		"SSH",
		"Owner",
	}
	var values [][]string

	for _, rp := range repos {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}

		values = append(
			values,
			[]string{
				rp.FullName,
				mode,
				rp.SSHURL,
				rp.Owner.UserName,
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}

func getRepoByPath(c *gitea.Client, repoPath string) (*gitea.Repository, error) {
	path := strings.Split(strings.Trim(repoPath, "/"), "/")
	switch len(path) {
	case 1:
		u, err := c.GetMyUserInfo()
		if err != nil {
			return nil, err
		}
		return c.GetRepo(u.UserName, path[0])
	case 2:
		return c.GetRepo(path[0], path[1])
	default:
		return nil, errors.New("repo path incorrect")
	}
}

func runRepoDetail(_ *cli.Context, path string) error {
	login := initCommandLoginOnly()
	client := login.Client()
	repo, err := getRepoByPath(client, path)
	if err != nil {
		return err
	}
	topics, err := client.ListRepoTopics(repo.Owner.UserName, repo.Name, gitea.ListRepoTopicsOptions{})
	if err != nil {
		return err
	}

	output := repo.FullName
	if repo.Mirror {
		output += " (mirror)"
	}
	if repo.Fork {
		output += " (fork)"
	}
	if repo.Archived {
		output += " (archived)"
	}
	if repo.Empty {
		output += " (empty)"
	}
	output += "\n"
	if len(topics) != 0 {
		output += "Topics: " + strings.Join(topics, ", ") + "\n"
	}
	output += "\n"
	output += repo.Description + "\n\n"
	output += fmt.Sprintf(
		"Open Issues: %d, Stars: %d, Forks: %d, Size: %s\n\n",
		repo.OpenIssues,
		repo.Stars,
		repo.Forks,
		utils.FormatSize(int64(repo.Size)),
	)

	fmt.Print(output)
	return nil
}

func runRepoCreate(ctx *cli.Context) error {
	login := initCommandLoginOnly()
	client := login.Client()
	var (
		repo *gitea.Repository
		err  error
	)
	opts := gitea.CreateRepoOption{
		Name:          ctx.String("name"),
		Description:   ctx.String("description"),
		Private:       ctx.Bool("private"),
		AutoInit:      ctx.Bool("init"),
		IssueLabels:   ctx.String("labels"),
		Gitignores:    ctx.String("gitignores"),
		License:       ctx.String("license"),
		Readme:        ctx.String("readme"),
		DefaultBranch: ctx.String("branch"),
	}
	if len(ctx.String("owner")) != 0 {
		repo, err = client.CreateOrgRepo(ctx.String("owner"), opts)
	} else {
		repo, err = client.CreateRepo(opts)
	}
	if err != nil {
		return err
	}
	if err = runRepoDetail(ctx, repo.FullName); err != nil {
		return err
	}
	fmt.Printf("%s\n", repo.HTMLURL)
	return nil
}
