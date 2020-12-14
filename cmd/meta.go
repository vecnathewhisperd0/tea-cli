// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/meta"
	"github.com/urfave/cli/v2"
)

// CmdMeta represents a sub command
var CmdMeta = cli.Command{
	Name:        "meta",
	Usage:       "Operations on tea itself",
	Description: `Operations on tea itself`,
	Subcommands: []*cli.Command{
		&meta.CmdAutocomplete,
	},
}
