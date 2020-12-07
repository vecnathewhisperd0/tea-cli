// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"
)

// CreatePull interactively creates a PR
func CreatePull(login *config.Login, ownerHint, repoHint string) error {
	var owner, repo, base, head, title, description string

	return task.CreatePull(
		login,
		owner,
		repo,
		base,
		head,
		title,
		description)
}
