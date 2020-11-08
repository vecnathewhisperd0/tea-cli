package task

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"

	"github.com/go-git/go-git/v5"
)

func PullCheckout(login *config.Login, repoOwner, repoName string, index int64) error {
	client := login.Client()

	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		return err
	}

	localBranchName, remoteBranchName, newRemoteName, remoteURL, err :=
		gitConfigForPR(localRepo, client, repoOwner, repoName, index)
	if err != nil {
		return err
	}

	// verify related remote is in local repo, otherwise add it
	localRemote, err := localRepo.GetOrCreateRemote(remoteURL, newRemoteName)
	if err != nil {
		return err
	}
	localRemoteName := localRemote.Config().Name

	// get auth & fetch remote
	fmt.Printf("Fetching PR %v (head %s:%s) from remote '%s'\n",
		index, remoteURL, remoteBranchName, localRemoteName)
	url, err := local_git.ParseURL(remoteURL)
	if err != nil {
		return err
	}
	auth, err := local_git.GetAuthForURL(url, login.User, login.SSHKey)
	if err != nil {
		return err
	}
	err = localRemote.Fetch(&git.FetchOptions{Auth: auth})
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println(err)
	} else if err != nil {
		return err
	}

	// checkout local branch
	fmt.Printf("Creating branch '%s'\n", localBranchName)
	err = localRepo.TeaCreateBranch(localBranchName, remoteBranchName, localRemoteName)
	if err == git.ErrBranchExists {
		fmt.Println("There may be changes since you last checked out, run `git pull` to get them.")
	} else if err != nil {
		return err
	}

	return localRepo.TeaCheckout(localBranchName)

	return nil
}

func gitConfigForPR(repo *local_git.TeaRepo, client *gitea.Client, owner, repoName string, idx int64) (localBranch, remoteBranch, remoteName, remoteURL string, err error) {
	// fetch PR source-repo & -branch from gitea
	pr, _, err := client.GetPullRequest(owner, repoName, idx)
	if err != nil {
		return
	}

	// test if we can pull via SSH, and configure git remote accordingly
	remoteURL = pr.Head.Repository.CloneURL
	keys, _, err := client.ListMyPublicKeys(gitea.ListPublicKeysOptions{})
	if err != nil {
		return
	}
	if len(keys) != 0 {
		remoteURL = pr.Head.Repository.SSHURL
	}

	// try to find a matching existing branch, otherwise return branch in pulls/ namespace
	localBranch = fmt.Sprintf("pulls/%v-%v", idx, pr.Head.Ref)
	if b, _ := repo.TeaFindBranchBySha(pr.Head.Sha, remoteURL); b != nil {
		localBranch = b.Name
	}

	remoteBranch = pr.Head.Ref
	remoteName = fmt.Sprintf("pulls/%v", pr.Head.Repository.Owner.UserName)
	return
}
