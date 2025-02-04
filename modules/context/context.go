// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package context

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"

	gogit "github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

var (
	errNotAGiteaRepo = errors.New("No Gitea login found. You might want to specify --repo (and --login) to work outside of a repository")
)

// TeaContext contains all context derived during command initialization and wraps cli.Context
type TeaContext struct {
	*cli.Context
	Login     *config.Login // config data & client for selected login
	RepoSlug  string        // <owner>/<repo>, optional
	Owner     string        // repo owner as derived from context or provided in flag, optional
	Repo      string        // repo name as derived from context or provided in flag, optional
	Output    string        // value of output flag
	LocalRepo *git.TeaRepo  // is set if flags specified a local repo via --repo, or if $PWD is a git repo
}

// GetListOptions return ListOptions based on PaginationFlags
func (ctx *TeaContext) GetListOptions() gitea.ListOptions {
	page := ctx.Int("page")
	limit := ctx.Int("limit")
	if limit < 0 {
		limit = 0
	}
	if limit != 0 && page == 0 {
		page = 1
	}
	return gitea.ListOptions{
		Page:     page,
		PageSize: limit,
	}
}

// GetRemoteRepoHTMLURL returns the web-ui url of the remote repo,
// after ensuring a remote repo is present in the context.
func (ctx *TeaContext) GetRemoteRepoHTMLURL() string {
	ctx.Ensure(CtxRequirement{RemoteRepo: true})
	return path.Join(ctx.Login.URL, ctx.Owner, ctx.Repo)
}

// Ensure checks if requirements on the context are set, and terminates otherwise.
func (ctx *TeaContext) Ensure(req CtxRequirement) {
	if req.LocalRepo && ctx.LocalRepo == nil {
		fmt.Println("Local repository required: Execute from a repo dir, or specify a path with --repo.")
		os.Exit(1)
	}

	if req.RemoteRepo && len(ctx.RepoSlug) == 0 {
		fmt.Println("Remote repository required: Specify ID via --repo or execute from a local git repo.")
		os.Exit(1)
	}
}

// CtxRequirement specifies context needed for operation
type CtxRequirement struct {
	// ensures a local git repo is available & ctx.LocalRepo is set. Implies .RemoteRepo
	LocalRepo bool
	// ensures ctx.RepoSlug, .Owner, .Repo are set
	RemoteRepo bool
}

// InitCommand resolves the application context, and returns the active login, and if
// available the repo slug. It does this by reading the config file for logins, parsing
// the remotes of the .git repo specified in repoFlag or $PWD, and using overrides from
// command flags. If a local git repo can't be found, repo slug values are unset.
func InitCommand(ctx *cli.Context) *TeaContext {
	// these flags are used as overrides to the context detection via local git repo
	repoFlag := ctx.String("repo")
	loginFlag := ctx.String("login")
	remoteFlag := ctx.String("remote")

	var (
		c                  TeaContext
		err                error
		repoPath           string // empty means PWD
		repoFlagPathExists bool
	)

	// check if repoFlag can be interpreted as path to local repo.
	if len(repoFlag) != 0 {
		if repoFlagPathExists, err = utils.DirExists(repoFlag); err != nil {
			log.Fatal(err.Error())
		}
		if repoFlagPathExists {
			repoPath = repoFlag
		}
	}

	if len(remoteFlag) == 0 {
		remoteFlag = config.GetPreferences().FlagDefaults.Remote
	}

	// try to read local git repo & extract context: if repoFlag specifies a valid path, read repo in that dir,
	// otherwise attempt PWD. if no repo is found, continue with default login
	if c.LocalRepo, c.Login, c.RepoSlug, err = contextFromLocalRepo(repoPath, remoteFlag); err != nil {
		if err == errNotAGiteaRepo || err == gogit.ErrRepositoryNotExists {
			// we can deal with that, commands needing the optional values use ctx.Ensure()
		} else {
			log.Fatal(err.Error())
		}
	}

	if len(repoFlag) != 0 && !repoFlagPathExists {
		// if repoFlag is not a valid path, use it to override repoSlug
		c.RepoSlug = repoFlag
	}

	// override login from flag, or use default login if repo based detection failed
	if len(loginFlag) != 0 {
		c.Login = config.GetLoginByName(loginFlag)
		if c.Login == nil {
			log.Fatalf("Login name '%s' does not exist", loginFlag)
		}
	} else if c.Login == nil {
		if c.Login, err = config.GetDefaultLogin(); err != nil {
			if err.Error() == "No available login" {
				// TODO: maybe we can directly start interact.CreateLogin() (only if
				// we're sure we can interactively!), as gh cli does.
				fmt.Println(`No gitea login configured. To start using tea, first run
  tea login add
and then run your command again.`)
			}
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "NOTE: no gitea login detected, falling back to login '%s'\n", c.Login.Name)
	}

	// parse reposlug (owner falling back to login owner if reposlug contains only repo name)
	c.Owner, c.Repo = utils.GetOwnerAndRepo(c.RepoSlug, c.Login.User)

	c.Context = ctx
	c.Output = ctx.String("output")
	return &c
}

// contextFromLocalRepo discovers login & repo slug from the default branch remote of the given local repo
func contextFromLocalRepo(repoPath, remoteValue string) (*git.TeaRepo, *config.Login, string, error) {
	repo, err := git.RepoFromPath(repoPath)
	if err != nil {
		return nil, nil, "", err
	}
	gitConfig, err := repo.Config()
	if err != nil {
		return repo, nil, "", err
	}

	if len(gitConfig.Remotes) == 0 {
		return repo, nil, "", errNotAGiteaRepo
	}

	// When no preferred value is given, choose a remote to find a
	// matching login based on its URL.
	if len(gitConfig.Remotes) > 1 && len(remoteValue) == 0 {
		// if master branch is present, use it as the default remote
		mainBranches := []string{"main", "master", "trunk"}
		for _, b := range mainBranches {
			masterBranch, ok := gitConfig.Branches[b]
			if ok {
				if len(masterBranch.Remote) > 0 {
					remoteValue = masterBranch.Remote
				}
				break
			}
		}
		// if no branch has matched, default to origin or upstream remote.
		if len(remoteValue) == 0 {
			if _, ok := gitConfig.Remotes["upstream"]; ok {
				remoteValue = "upstream"
			} else if _, ok := gitConfig.Remotes["origin"]; ok {
				remoteValue = "origin"
			}
		}
	}
	// make sure a remote is selected
	if len(remoteValue) == 0 {
		for remote := range gitConfig.Remotes {
			remoteValue = remote
			break
		}
	}

	remoteConfig, ok := gitConfig.Remotes[remoteValue]
	if !ok || remoteConfig == nil {
		return repo, nil, "", fmt.Errorf("Remote '%s' not found in this Git repository", remoteValue)
	}

	logins, err := config.GetLogins()
	if err != nil {
		return repo, nil, "", err
	}
	for _, l := range logins {
		sshHost := l.GetSSHHost()
		for _, u := range remoteConfig.URLs {
			p, err := git.ParseURL(u)
			if err != nil {
				return repo, nil, "", fmt.Errorf("Git remote URL parse failed: %s", err.Error())
			}
			if strings.EqualFold(p.Scheme, "http") || strings.EqualFold(p.Scheme, "https") {
				if strings.HasPrefix(u, l.URL) {
					ps := strings.Split(p.Path, "/")
					path := strings.Join(ps[len(ps)-2:], "/")
					return repo, &l, strings.TrimSuffix(path, ".git"), nil
				}
			} else if strings.EqualFold(p.Scheme, "ssh") {
				if sshHost == p.Host {
					return repo, &l, strings.TrimLeft(p.Path, "/"), nil
				}
			}
		}
	}

	return repo, nil, "", errNotAGiteaRepo
}
