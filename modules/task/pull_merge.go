// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
)

// PullMerge merges a PR
func PullMerge(login *config.Login, repoOwner, repoName string, index int64, opt gitea.MergePullRequestOption) error {
	client := login.Client()
	success, _, err := client.MergePullRequest(repoOwner, repoName, index, opt)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("Failed to merge PR. Is it still open?")
	}
	return nil
}
