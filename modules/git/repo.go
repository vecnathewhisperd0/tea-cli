// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// TeaRepo is a go-git Repository, with an extended high level interface.
type TeaRepo struct {
	*git.Repository
}

// RepoForWorkdir tries to open the git repository in the local directory
// for reading or modification.
func RepoForWorkdir() (*TeaRepo, error) {
	return RepoFromPath("")
}

// RepoFromPath tries to open the git repository by path
func RepoFromPath(path string) (*TeaRepo, error) {
	if len(path) == 0 {
		path = "./"
	}
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	return &TeaRepo{repo}, nil
}

// CloneIntoPath emulates a git clone <url> --depth <n> ...
func CloneIntoPath(url, path string, auth transport.AuthMethod, depth int, insecure bool) (*TeaRepo, error) {
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:             url,
		Auth:            auth,
		Depth:           depth,
		InsecureSkipTLS: insecure,
	})
	if err != nil {
		return nil, err
	}
	return &TeaRepo{repo}, nil
}
