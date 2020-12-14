// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// editPullState abstracts the arg parsing to edit the given pull request
func editPullState(ctx *cli.Context, opts gitea.EditPullRequestOption) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	if ctx.Args().Len() == 0 {
		log.Fatal(ctx.Command.ArgsUsage)
	}

	index, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	pr, _, err := login.Client().EditPullRequest(owner, repo, index, opts)
	if err != nil {
		return err
	}

	print.PullDetails(pr, nil)
	return nil
}
