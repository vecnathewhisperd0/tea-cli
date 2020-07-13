// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	git_config "github.com/go-git/go-git/v5/config"
	git_plumbing "github.com/go-git/go-git/v5/plumbing"
	git_transport "github.com/go-git/go-git/v5/plumbing/transport"
)

// TeaCreateBranch creates a new branch in the repo, tracking from another branch.
func (r TeaRepo) TeaCreateBranch(localBranchName, remoteBranchName, remoteName string) error {
	// save in .git/config to assign remote for future pulls
	localBranchRefName := git_plumbing.NewBranchReferenceName(localBranchName)
	err := r.CreateBranch(&git_config.Branch{
		Name:   localBranchName,
		Merge:  git_plumbing.NewBranchReferenceName(remoteBranchName),
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
func (r TeaRepo) TeaDeleteBranch(branch *git_config.Branch, remoteBranch string, auth git_transport.AuthMethod) error {
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
			Auth:       auth,
		})
	}

	return err
}

// TeaFindBranchBySha returns a branch that is at the the given SHA and syncs to the
// given remote repo.
func (r TeaRepo) TeaFindBranchBySha(sha, repoURL string) (b *git_config.Branch, err error) {
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
			if ref.Hash().String() == sha && strings.HasPrefix(name, remoteName) {
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
	if remoteRefName == "" || localRefName == "" {
		// no remote tracking branch found, so a potential local branch
		// can't be a match either
		return nil, nil
	}

	b = &git_config.Branch{
		Remote: remoteName,
		Name:   localRefName.Short(),
		Merge:  localRefName,
	}
	return b, b.Validate()
}

// TeaFindBranchByName returns a branch that is at the the given local and
// remote names and syncs to the given remote repo. This method is less precise
// than TeaFindBranchBySha(), but may be desirable if local and remote branch
// have diverged.
func (r TeaRepo) TeaFindBranchByName(branchName, repoURL string) (b *git_config.Branch, err error) {
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
	var remoteSearchingName = fmt.Sprintf("%s/%s", remoteName, branchName)
	err = iter.ForEach(func(ref *git_plumbing.Reference) error {
		if ref.Name().IsRemote() && ref.Name().Short() == remoteSearchingName {
			remoteRefName = ref.Name()
		}
		n := ref.Name()
		if n.IsBranch() && n.Short() == branchName {
			localRefName = n
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if remoteRefName == "" || localRefName == "" {
		return nil, nil
	}

	b = &git_config.Branch{
		Remote: remoteName,
		Name:   localRefName.Short(),
		Merge:  localRefName,
	}
	return b, b.Validate()
}

// TeaGetCurrentBranchName return the name of the branch witch is currently active
func (r TeaRepo) TeaGetCurrentBranchName() string {
	localHead, err := r.Head()
	if err != nil {
		return ""
	}
	if localHead.Type() != git_plumbing.SymbolicReference || !localHead.Name().IsBranch() {
		return ""
	}

	return localHead.Name().String()
}
