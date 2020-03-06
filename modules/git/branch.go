package git

import (
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
