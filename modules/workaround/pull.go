// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package workaround

import (
	"net/url"

	"code.gitea.io/sdk/gitea"
)

// FixPullHeadSha is a workaround for https://github.com/go-gitea/gitea/issues/12675
func FixPullHeadSha(client *gitea.Client, pr *gitea.PullRequest, repoOwner, repoName string) error {
	if pr.Head != nil && pr.Head.Sha == "" {
		headCommit, resp, err := client.GetSingleCommit(repoOwner, repoName, url.PathEscape(pr.Head.Ref))
		if err != nil && resp == nil || err != nil && resp.StatusCode != 404 {
			return err
		}
		if headCommit != nil {
			pr.Head.Sha = headCommit.SHA
		}
	}
	return nil
}
