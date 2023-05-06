// Copyright 2023 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"
	"github.com/AlecAivazis/survey/v2"
)

// MergePull interactively creates a PR
func MergePull(ctx *context.TeaContext) error {
	if ctx.LocalRepo == nil {
		return fmt.Errorf("Must specify a PR index")
	}

	branch, _, err := ctx.LocalRepo.TeaGetCurrentBranchNameAndSHA()
	if err != nil {
		return err
	}

	idx, err := getPullIndex(ctx, branch)
	if err != nil {
		return err
	}

	return task.PullMerge(ctx.Login, ctx.Owner, ctx.Repo, idx, gitea.MergePullRequestOption{
		Style:   gitea.MergeStyle(ctx.String("style")),
		Title:   ctx.String("title"),
		Message: ctx.String("message"),
	})
}

// getPullIndex interactively determines the PR index
func getPullIndex(ctx *context.TeaContext, branch string) (int64, error) {
	prs, _, err := ctx.Login.Client().ListRepoPullRequests(ctx.Owner, ctx.Repo, gitea.ListPullRequestsOptions{
		State: gitea.StateOpen,
	})
	if err != nil {
		return 0, err
	}
	if len(prs) == 0 {
		return 0, fmt.Errorf("No open PRs found")
	}

	prOptions := make([]string, 0)

	// get the PR indexes where head branch is the current branch
	for _, pr := range prs {
		if pr.Head.Ref == branch {
			prOptions = append(prOptions, fmt.Sprintf("#%d: %s", pr.Index, pr.Title))
		}
	}

	// then get the PR indexes where base branch is the current branch
	for _, pr := range prs {
		// don't add the same PR twice, so `pr.Head.Ref != branch`
		if pr.Base.Ref == branch && pr.Head.Ref != branch {
			prOptions = append(prOptions, fmt.Sprintf("#%d: %s", pr.Index, pr.Title))
		}
	}

	selected := ""
	q := &survey.Select{
		Message:  "Select a PR to merge",
		Options:  prOptions,
		PageSize: 10,
	}
	err = survey.AskOne(q, &selected)
	if err != nil {
		return 0, err
	}

	// get the index from the selected option
	before, _, _ := strings.Cut(selected, ":")
	before = strings.TrimPrefix(before, "#")
	idx, err := utils.ArgToIndex(before)
	if err != nil {
		return 0, err
	}

	return idx, nil
}
