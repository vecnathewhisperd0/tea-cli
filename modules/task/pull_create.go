// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"
	"log"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"github.com/go-git/go-git/v5"
)

// PullCreate creates a PR in the given repo and prints the result
func CreatePull(login *config.Login, repoOwner, repoName, base, head, title, description string) error {
	client := login.Client()

	repo, _, err := client.GetRepo(repoOwner, repoName)
	if err != nil {
		log.Fatal("could not fetch repo meta: ", err)
	}

	// open local git repo
	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		log.Fatal("could not open local repo: ", err)
	}

	// push if possible
	log.Println("git push")
	err = localRepo.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Printf("Error occurred during 'git push':\n%s\n", err.Error())
	}

	// default is default branch
	if len(base) == 0 {
		base = repo.DefaultBranch
	}

	// default is current one
	if len(head) == 0 {
		headBranch, err := localRepo.Head()
		if err != nil {
			log.Fatal(err)
		}
		sha := headBranch.Hash().String()

		remote, err := localRepo.TeaFindBranchRemote("", sha)
		if err != nil {
			log.Fatal("could not determine remote for current branch: ", err)
		}

		if remote == nil {
			// if no remote branch is found for the local hash, we abort:
			// user has probably not configured a remote for the local branch,
			// or local branch does not represent remote state.
			log.Fatal("no matching remote found for this branch. try git push -u <remote> <branch>")
		}

		branchName, err := localRepo.TeaGetCurrentBranchName()
		if err != nil {
			log.Fatal(err)
		}

		url, err := local_git.ParseURL(remote.Config().URLs[0])
		if err != nil {
			log.Fatal(err)
		}
		owner, _ := utils.GetOwnerAndRepo(strings.TrimLeft(url.Path, "/"), "")
		if owner != repo.Owner.UserName {
			head = fmt.Sprintf("%s:%s", owner, branchName)
		} else {
			head = branchName
		}
	}

	// head & base may not be the same
	if head == base {
		return fmt.Errorf("Can't create PR from %s to %s\n", head, base)
	}

	// default is head branch name
	if len(title) == 0 {
		title = head
		if strings.Contains(title, ":") {
			title = strings.SplitN(title, ":", 2)[1]
		}
		title = strings.Replace(title, "-", " ", -1)
		title = strings.Replace(title, "_", " ", -1)
		title = strings.Title(strings.ToLower(title))
	}
	// title is required
	if len(title) == 0 {
		return fmt.Errorf("Title is required")
	}

	pr, _, err := client.CreatePullRequest(repoOwner, repoName, gitea.CreatePullRequestOption{
		Head:  head,
		Base:  base,
		Title: title,
		Body:  description,
	})

	if err != nil {
		log.Fatalf("could not create PR from %s to %s:%s: %s", head, repoOwner, base, err)
	}

	print.PullDetails(pr)

	fmt.Println(pr.HTMLURL)

	return err
}
