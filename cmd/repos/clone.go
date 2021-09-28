// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/interact"

	git_config "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/urfave/cli/v2"
)

// CmdRepoClone represents a sub command of repos to create a local copy
var CmdRepoClone = cli.Command{
	Name:        "clone",
	Aliases:     []string{"C"},
	Usage:       "Clone a repository locally",
	Description: "Clone a repository locally, without a local git installation required (defaults to PWD)",
	Action:      runRepoClone,
	ArgsUsage:   "[target dir]",
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:    "depth",
			Aliases: []string{"d"},
			Usage:   "num commits to fetch, defaults to all",
		},
	}, flags.LoginRepoFlags...),
}

func runRepoClone(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	repoMeta, _, err := ctx.Login.Client().GetRepo(ctx.Owner, ctx.Repo)
	if err != nil {
		return err
	}

	originURL, err := cloneURL(repoMeta, ctx.Login)
	if err != nil {
		return err
	}

	auth, err := git.GetAuthForURL(originURL, ctx.Login.Token, ctx.Login.SSHKey, interact.PromptPassword)
	if err != nil {
		return err
	}

	// default path behaviour as native git
	localPath := ctx.Args().First()
	if localPath == "" {
		localPath = ctx.Repo
	}

	repo, err := git.CloneIntoPath(
		originURL.String(),
		localPath,
		auth,
		ctx.Int("depth"),
		ctx.Login.Insecure,
	)

	// set up upstream remote for forks
	if repoMeta.Fork && repoMeta.Parent != nil {
		upstreamURL, err := cloneURL(repoMeta.Parent, ctx.Login)
		if err != nil {
			return err
		}
		upstreamBranch := repoMeta.Parent.DefaultBranch
		repo.CreateRemote(&git_config.RemoteConfig{
			Name: "upstream",
			URLs: []string{upstreamURL.String()},
		})
		repoConf, err := repo.Config()
		if err != nil {
			return err
		}
		if b, ok := repoConf.Branches[upstreamBranch]; ok {
			b.Remote = "upstream"
			b.Merge = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", upstreamBranch))
		}
		return repo.SetConfig(repoConf)
	}

	return err
}

func cloneURL(repo *gitea.Repository, login *config.Login) (*url.URL, error) {
	urlStr := repo.CloneURL
	if login.SSHKey != "" {
		urlStr = repo.SSHURL
	}
	return git.ParseURL(urlStr)
}
