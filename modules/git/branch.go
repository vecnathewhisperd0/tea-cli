// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"fmt"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	git_config "gopkg.in/src-d/go-git.v4/config"
	git_plumbing "gopkg.in/src-d/go-git.v4/plumbing"
)

// TeaCreateBranch creates a new branch in the repo, tracking from another branch.
// If remoteName is not-null, a remote branch is tracked.
func (r TeaRepo) TeaCreateBranch(localBranchName, remoteBranchName, remoteName string) error {
	// save in .git/config to assign remote for future pulls
	localBranchRefName := git_plumbing.NewBranchReferenceName(localBranchName)
	err := r.CreateBranch(&git_config.Branch{
		Name:   localBranchName,
		Merge:  git_plumbing.NewBranchReferenceName(remoteBranchName), // FIXME: should be remoteBranchName
		Remote: remoteName,
	})
	if err != nil {
		return err
	}

	// serialize the branch to .git/refs/heads
	remoteBranchRefName := git_plumbing.NewRemoteReferenceName(remoteName, remoteBranchName)
	remoteBranchRef, err := r.Storer.Reference(remoteBranchRefName)
	if err != nil {
		return err
	}
	localHashRef := git_plumbing.NewHashReference(localBranchRefName, remoteBranchRef.Hash())
	return r.Storer.SetReference(localHashRef)
}

// TeaCheckout checks out the given branch in the worktree.
func (r TeaRepo) TeaCheckout(branchName string) error {
	tree, err := r.Worktree()
	if err != nil {
		return err
	}
	localBranchRefName := git_plumbing.NewBranchReferenceName(branchName)
	return tree.Checkout(&git.CheckoutOptions{Branch: localBranchRefName})
}

// TeaDeleteBranch removes the given branch locally, and if `remoteBranch` is
// not empty deletes it at it's remote repo.
func (r TeaRepo) TeaDeleteBranch(branch *git_config.Branch, remoteBranch string) error {
	err := r.DeleteBranch(branch.Name)
	// if the branch is not found that's ok, as .git/config may have no entry if
	// no remote tracking branch is configured for it (eg push without -u flag)
	if err != nil && err.Error() != "branch not found" {
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
	// find remote matching our repoURL
	remote, err := r.GetRemote(repoURL)
	if err != nil {
		return nil, err
	}
	if remote == nil {
		return nil, fmt.Errorf("No remote found for '%s'", repoURL)
	}
	remoteName := remote.Config().Name

	// check if the given remote has our branch (.git/refs/remotes/<remoteName>/*)
	iter, err := r.References()
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	var remoteRefName git_plumbing.ReferenceName
	var localRefName git_plumbing.ReferenceName
	err = iter.ForEach(func(ref *git_plumbing.Reference) error {
		if ref.Name().IsRemote() {
			name := ref.Name().Short()
			if name != "master" &&
				ref.Hash().String() == sha &&
				strings.HasPrefix(name, remoteName) {
				remoteRefName = ref.Name()
			}
		}

		if ref.Name().IsBranch() && ref.Hash().String() == sha {
			localRefName = ref.Name()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if remoteRefName == "" {
		// no remote tracking branch found, so a potential local branch
		// can't be a match either
		return nil, nil
	}

	b = &git_config.Branch{
		Remote: remoteName,
		Name:   localRefName.Short(),
		Merge:  localRefName,
	}
	fmt.Println(b)
	return b, b.Validate()
}
