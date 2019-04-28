// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/sdk/gitea"
	"github.com/spf13/cobra"
)

var loginName string
var repoPath string

func init() {
	rootCmd.AddCommand(pullCmd)
	rootCmd.PersistentFlags().StringVarP(&loginName, "login", "l", "", "Indicate one login, optional when inside a gitea repository")
	rootCmd.PersistentFlags().StringVarP(&repoPath, "repo", "r", "", "Indicate one repository, optional when inside a gitea repository")
}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pulls",
	Short: "Operate with pulls of the repository",
	Long:  `Operate with pulls of the repository`,
	Run:   runPulls,
}

func runPulls(cmd *cobra.Command, args []string) {
	login := getLoginByName(loginName)
	if login == nil {
		Errorf("Login '%s' not found in config\n", loginName)
		return
	}
	owner, repo := getOwnerRepo()

	prs, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		Page:  0,
		State: string(gitea.StateOpen),
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(prs) == 0 {
		fmt.Println("No pull requests left")
	}

	for _, pr := range prs {
		if pr == nil {
			continue
		}
		name := pr.Poster.FullName
		if len(name) == 0 {
			name = pr.Poster.UserName
		}
		fmt.Printf("#%d\t%s\t%s\t%s\n", pr.Index, name, pr.Updated.Format("2006-01-02 15:04:05"), pr.Title)
	}
}
