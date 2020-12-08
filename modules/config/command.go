// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"
)

// InitCommand returns repository and *Login based on flags
func InitCommand(repoValue, loginValue, remoteValue string) (*Login, string, string) {
	var login *Login

	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if login, err = GetDefaultLogin(); err != nil {
		log.Fatal(err.Error())
	}

	exist, err := utils.PathExists(repoValue)
	if err != nil {
		log.Fatal(err.Error())
	}

	if exist || len(repoValue) == 0 {
		login, repoValue, err = curGitRepoPath(repoValue, remoteValue)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if loginValue != "" {
		login = GetLoginByName(loginValue)
		if login == nil {
			log.Fatal("Login name " + loginValue + " does not exist")
		}
	}

	owner, repo := utils.GetOwnerAndRepo(repoValue, login.User)
	return login, owner, repo
}

// InitCommandLoginOnly return *Login based on flags
func InitCommandLoginOnly(loginValue string) *Login {
	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var login *Login
	if loginValue == "" {
		login, err = GetDefaultLogin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		login = GetLoginByName(loginValue)
		if login == nil {
			log.Fatal("Login name " + loginValue + " does not exist")
		}
	}

	return login
}

func curGitRepoPath(repoValue, remoteValue string) (*Login, string, error) {
	var err error
	var repo *git.TeaRepo
	if len(repoValue) == 0 {
		repo, err = git.RepoForWorkdir()
	} else {
		repo, err = git.RepoFromPath(repoValue)
	}
	if err != nil {
		return nil, "", err
	}
	gitConfig, err := repo.Config()
	if err != nil {
		return nil, "", err
	}

	// if no remote
	if len(gitConfig.Remotes) == 0 {
		return nil, "", errors.New("No remote(s) found in this Git repository")
	}

	// if only one remote exists
	if len(gitConfig.Remotes) >= 1 && len(remoteValue) == 0 {
		for remote := range gitConfig.Remotes {
			remoteValue = remote
		}
		if len(gitConfig.Remotes) > 1 {
			// if master branch is present, use it as the default remote
			masterBranch, ok := gitConfig.Branches["master"]
			if ok {
				if len(masterBranch.Remote) > 0 {
					remoteValue = masterBranch.Remote
				}
			}
		}
	}

	remoteConfig, ok := gitConfig.Remotes[remoteValue]
	if !ok || remoteConfig == nil {
		return nil, "", errors.New("Remote " + remoteValue + " not found in this Git repository")
	}

	for _, l := range Config.Logins {
		for _, u := range remoteConfig.URLs {
			p, err := git.ParseURL(strings.TrimSpace(u))
			if err != nil {
				return nil, "", fmt.Errorf("Git remote URL parse failed: %s", err.Error())
			}
			if strings.EqualFold(p.Scheme, "http") || strings.EqualFold(p.Scheme, "https") {
				if strings.HasPrefix(u, l.URL) {
					ps := strings.Split(p.Path, "/")
					path := strings.Join(ps[len(ps)-2:], "/")
					return &l, strings.TrimSuffix(path, ".git"), nil
				}
			} else if strings.EqualFold(p.Scheme, "ssh") {
				if l.GetSSHHost() == strings.Split(p.Host, ":")[0] {
					return &l, strings.TrimLeft(strings.TrimSuffix(p.Path, ".git"), "/"), nil
				}
			}
		}
	}

	return nil, "", errors.New("No Gitea login found. You might want to specify --repo (and --login) to work outside of a repository")
}
