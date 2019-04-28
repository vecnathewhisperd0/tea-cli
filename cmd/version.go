// Copyright © 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Version holds the current Gitea version
var Version = "0.1.0-dev"

// Tags holds the build tags used
var Tags = ""

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(Tags) > 0 {
			Version += " built with: " + strings.Replace(Tags, " ", ", ", -1)
		}
		fmt.Println("Version " + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
