// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	git_config "gopkg.in/src-d/go-git.v4/config"
	git_plumbing "gopkg.in/src-d/go-git.v4/plumbing"
)

// TeaCreateBranch creates a new branch in the repo, tracking from another branch.
// If remoteName is not-null, a remote branch is tracked.
func (r TeaRepo) TeaCreateBranch(localBranchName, remoteBranchName, remoteName string) error {
	remoteBranchRefName := git_plumbing.NewBranchReferenceName(remoteBranchName)
	err := r.CreateBranch(&git_config.Branch{
		Name:   localBranchName,
		Merge:  remoteBranchRefName,
		Remote: remoteName,
	})
	if err != nil {
		return err
	}

	// serialize the branch to .git/refs/heads (otherwise branch is only defined
	// in .git/.config)
	localBranchRefName := git_plumbing.NewBranchReferenceName(localBranchName)
	remoteBranchRef, err := r.Storer.Reference(remoteBranchRefName)
	if err != nil {
		return nil
	}
	localHashRef := git_plumbing.NewHashReference(localBranchRefName, remoteBranchRef.Hash())
	r.Storer.SetReference(localHashRef)
	return nil
}

// TeaCheckout checks out the given branch in the worktree.
func (r TeaRepo) TeaCheckout(branchName string) error {
	tree, err := r.Worktree()
	if err != nil {
		return nil
	}
	localBranchRefName := git_plumbing.NewBranchReferenceName(branchName)
	return tree.Checkout(&git.CheckoutOptions{Branch: localBranchRefName})
}

// TeaDeleteBranch removes the given branch locally, and if `remoteBranch` is
// not empty deletes it at it's remote repo.
func (r TeaRepo) TeaDeleteBranch(branch *git_config.Branch, remoteBranch string) error {
	err := r.DeleteBranch(branch.Name)
	if err != nil {
		return err
	}
	err = r.Storer.RemoveReference(git_plumbing.NewBranchReferenceName(branch.Name))
	if err != nil {
		return err
	}

	if remoteBranch != "" {
		// delete remote branch via git protocol:
		// an empty source in the refspec means remote deletion to git 🙃
		refspec := fmt.Sprintf(":%s", git_plumbing.NewBranchReferenceName(remoteBranch))
		err = r.Push(&git.PushOptions{
			RemoteName: branch.Remote,
			RefSpecs:   []git_config.RefSpec{git_config.RefSpec(refspec)},
			Prune:      true,
		})
	}

	return err
}

// TeaFindBranch returns a branch that is at the the given SHA and syncs to the
// given remote repo.
func (r TeaRepo) TeaFindBranch(sha, repoURL string) (b *git_config.Branch, err error) {
	url, err := ParseURL(repoURL)
	if err != nil {
		return nil, err
	}

	branches, err := r.Branches()
	if err != nil {
		return nil, err
	}
	err = branches.ForEach(func(ref *git_plumbing.Reference) error {
		name := ref.Name().Short()
		if name != "master" && ref.Hash().String() == sha {
			branch, _ := r.Branch(name)
			repoConf, err := r.Config()
			if err != nil {
				return err
			}
			remote := repoConf.Remotes[branch.Remote]
			for _, u := range remote.URLs {
				remoteURL, err := ParseURL(u)
				if err != nil {
					return err
				}
				if remoteURL.Host == url.Host && remoteURL.Path == url.Path {
					b = branch
				}
			}
		}
		return nil
	})

	return b, err
}
