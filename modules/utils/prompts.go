// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func PromptRepoSlug(defaultOwner, defaultRepo string) (owner, repo string, err error) {
	prompt := "Target repo:"
	required := true
	if len(defaultOwner) != 0 && len(defaultRepo) != 0 {
		prompt = fmt.Sprintf("Target repo [%s/%s]:", defaultOwner, defaultRepo)
		required = false
	}
	var repoSlug string

	owner = defaultOwner
	repo = defaultRepo

	err = survey.AskOne(
		&survey.Input{Message: prompt},
		&repoSlug,
		survey.WithValidator(func(input interface{}) error {
			if str, ok := input.(string); ok {
				if !required && len(str) == 0 {
					return nil
				}
				split := strings.Split(str, "/")
				if len(split) != 2 || len(split[0]) == 0 || len(split[1]) == 0 {
					return fmt.Errorf("must follow the <owner>/<repo> syntax")
				}
			} else {
				return fmt.Errorf("invalid result type")
			}
			return nil
		}),
	)

	if err == nil && len(repoSlug) != 0 {
		repoSlugSplit := strings.Split(repoSlug, "/")
		owner = repoSlugSplit[0]
		repo = repoSlugSplit[1]
	}
	return
}
