// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdPullsMerge merges a PR
var CmdPullsMerge = cli.Command{
	Name:        "merge",
	Aliases:     []string{"m"},
	Usage:       "Merge a pull request",
	Description: "Merge a pull request",
	ArgsUsage:   "<pull index>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "style",
			Aliases: []string{"s"},
			Usage:   "Kind of merge to perform: merge, rebase, squash, rebase-merge",
			Value:   "merge",
		},
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "Merge commit title",
		},
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
			Usage:   "Merge commit message",
		},
	}, flags.AllDefaultFlags...),
	Action: func(cmd *cli.Context) error {
		ctx := context.InitCommand(cmd)
		ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

		var idx int64
		var err error
		if ctx.Args().Len() == 1 {
			idx, err = utils.ArgToIndex(ctx.Args().First())
			if err != nil {
				return err
			}
		} else {
			if ctx.LocalRepo == nil {
				return fmt.Errorf("Must specify a PR index")
			}

			branch, _, err := ctx.LocalRepo.TeaGetCurrentBranchNameAndSHA()
			if err != nil {
				return err
			}

			idx, err = getPullIndexByBranch(ctx, branch)
			if err != nil {
				return err
			}
		}

		success, _, err := ctx.Login.Client().MergePullRequest(ctx.Owner, ctx.Repo, idx, gitea.MergePullRequestOption{
			Style:   gitea.MergeStyle(ctx.String("style")),
			Title:   ctx.String("title"),
			Message: ctx.String("message"),
		})

		if err != nil {
			return err
		}
		if !success {
			return fmt.Errorf("Failed to merge PR. Is it still open?")
		}
		return nil
	},
}

func getPullIndexByBranch(ctx *context.TeaContext, branch string) (int64, error) {
	prs, _, err := ctx.Login.Client().ListRepoPullRequests(ctx.Owner, ctx.Repo, gitea.ListPullRequestsOptions{
		State: gitea.StateOpen,
	})
	if err != nil {
		return 0, err
	}

	for _, pr := range prs {
		if pr.Head.Ref == branch {
			return pr.Index, nil
		}
	}
	return 0, fmt.Errorf("No open PR for branch %s", branch)
}
