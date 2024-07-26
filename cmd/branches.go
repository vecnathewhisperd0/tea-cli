// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"code.gitea.io/tea/cmd/branches"

	"github.com/urfave/cli/v2"
)

// CmdBranches represents to login a gitea server.
var CmdBranches = cli.Command{
	Name:        "branches",
	Aliases:     []string{"branch", "b"},
	Category:    catEntities,
	Usage:       "Consult branches",
	Description: `Lists branches when called without argument. If a branch is provided, will show it in detail.`,
	ArgsUsage:   "[<branch name>]",
	Action:      runBranches,
	Subcommands: []*cli.Command{
		&branches.CmdBranchesList,
		&branches.CmdBranchesProtect,
		&branches.CmdBranchesUnprotect,
	},
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:  "comments",
			Usage: "Whether to display comments (will prompt if not provided & run interactively)",
		},
	}, branches.CmdBranchesList.Flags...),
}

func runBranches(ctx *cli.Context) error {
	return branches.RunBranchesList(ctx)
}
