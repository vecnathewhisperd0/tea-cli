// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/modules/interact"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/utils"
	"github.com/urfave/cli/v2"
)

// CmdAddComment is the main command to operate with notifications
var CmdAddComment = cli.Command{
	Name:        "comment",
	Usage:       "Add a comment to an issue / pr",
	Description: "Add a comment to an issue / pr",
	ArgsUsage:   "<issue / pr index> [<comment body>]",
	Action:      runAddComment,
	Flags:       flags.AllDefaultFlags,
}

func runAddComment(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	args := ctx.Args()
	if args.Len() == 0 {
		return fmt.Errorf("Please specify issue / pr index")
	}

	idx, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	body := strings.Join(ctx.Args().Tail(), " ")
	if len(body) == 0 {
		if body, err = interact.PromptMultiline("Content"); err != nil {
			return err
		}
	}

	client := ctx.Login.Client()
	_, _, err = client.CreateIssueComment(ctx.Owner, ctx.Repo, idx, gitea.CreateIssueCommentOption{
		Body: body,
	})
	return err
}
