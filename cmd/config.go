// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/sdk/gitea"
	"github.com/muesli/termenv"
	"github.com/urfave/cli/v2"
)

func getGlamourTheme() string {
	if termenv.HasDarkBackground() {
		return "dark"
	}
	return "light"
}

func getListOptions(ctx *cli.Context) gitea.ListOptions {
	page := ctx.Int("page")
	limit := ctx.Int("limit")
	if limit != 0 && page == 0 {
		page = 1
	}
	return gitea.ListOptions{
		Page:     page,
		PageSize: limit,
	}
}
