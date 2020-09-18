// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"path"
	"strconv"
	"strings"

	local_git "code.gitea.io/tea/modules/git"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

// CmdOpen represents a sub command of issues to open issue on the web browser
var CmdOpen = cli.Command{
	Name:        "open",
	Usage:       "Open something of the repository on web browser",
	Description: `Open something of the repository on web browser`,
	Action:      runOpen,
	Flags:       append([]cli.Flag{}, LoginRepoFlags...),
	ArgsUsage:   "[<owner>/<repo>] [<issue index> | issues | pulls | releases | commits | branches | wiki | activity | settings | labels | milestones]",
}

func runOpen(ctx *cli.Context) error {
	login := initCommandLoginOnly()
	arg1 := ctx.Args().Get(0)
	arg2 := ctx.Args().Get(1)
	owner, repo, view, index := argsToIndexOrRepo(arg1, arg2, repoValue)

	// if no repo specified via args / flags, try to extract from .git in PWD
	if owner == "" {
		l, ownerAndRepo, err := curGitRepoPath(repoValue)
		if err == nil {
			owner, repo = getOwnerAndRepo(ownerAndRepo, login.User)
			login = l
		}
	}

	if owner != "" && repo != "" {
		switch {
		case strings.EqualFold(arg1, "issues"):
			view = "issues"
		case strings.EqualFold(arg1, "pulls"):
			view = "pulls"
		case strings.EqualFold(arg1, "releases"):
			view = "releases"
		case strings.EqualFold(arg1, "commits"):
			view = getCommitView()
		case strings.EqualFold(arg1, "branches"):
			view = "branches"
		case strings.EqualFold(arg1, "wiki"):
			view = "wiki"
		case strings.EqualFold(arg1, "activity"):
			view = "activity"
		case strings.EqualFold(arg1, "settings"):
			view = "settings"
		case strings.EqualFold(arg1, "labels"):
			view = "labels"
		case strings.EqualFold(arg1, "milestones"):
			view = "milestones"
		case index != "":
			view = "issues/" + index
		}
	} else if index != "" {
		log.Fatal("no repository specified")
	}

	u := path.Join(login.URL, owner, repo, view)
	err := open.Start(u)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// parses an argument either to an issue index, or a owner/repo combination
func argsToIndexOrRepo(arg1, arg2, overrideRepoSlug string) (owner, repo, view, idx string) {
	owner, repo = getOwnerAndRepo(overrideRepoSlug, "")

	_, err := strconv.ParseInt(arg1, 10, 64)
	arg1IsIndex := err == nil
	_, err = strconv.ParseInt(arg2, 10, 64)
	arg2IsIndex := err == nil

	if arg1IsIndex {
		idx = arg1
		return
	} else if repo != "" {
		view = arg1 // view was already provided by overrideRepoSlug, so we use arg1 as view param
	} else {
		owner, repo = getOwnerAndRepo(arg1, "")
	}

	if arg2IsIndex {
		idx = arg2
	} else {
		view = arg2
	}

	return
}

func getCommitView() string {
	repo, err := local_git.RepoForWorkdir()
	if err != nil {
		log.Fatal(err)
	}
	b, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	name := b.Name()
	switch {
	case name.IsBranch():
		return "commits/branch/" + name.Short()
	case name.IsTag():
		return "commits/tag/" + name.Short()
	default:
		return ""
	}
}
