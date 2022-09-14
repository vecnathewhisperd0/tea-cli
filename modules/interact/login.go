// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"regexp"
	"strings"

	"code.gitea.io/tea/modules/interact/prompts"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
)

// loginMethod does represent methods tea do support to auth against gitea
type loginMethod string

const (
	loginMethodToken    loginMethod = "access token"
	loginMethodSSH      loginMethod = "ssh-key or -certificate"
	loginMethodPassword loginMethod = "username & password"
)

// CreateLogin create an login interactive
func CreateLogin() error {
	var name, token, user, passwd, sshKey, giteaURL, sshCertPrincipal, sshKeyFingerprint string
	insecure := false
	sshAgent := false

	promptI := &survey.Input{Message: "URL of Gitea instance: "}
	if err := survey.AskOne(promptI, &giteaURL, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	giteaURL = strings.TrimSuffix(strings.TrimSpace(giteaURL), "/")
	if len(giteaURL) == 0 {
		fmt.Println("URL is required!")
		return nil
	}

	name, err := task.GenerateLoginName(giteaURL, "")
	if err != nil {
		return err
	}

	promptI = &survey.Input{Message: "Name of new Login [" + name + "]: "}
	if err := survey.AskOne(promptI, &name); err != nil {
		return err
	}

	lMethod, err := prompts.Select("Login with: ", []string{
		string(loginMethodSSH),
		string(loginMethodToken),
		string(loginMethodPassword),
	}, "", "")
	if err != nil {
		return err
	}

	switch loginMethod(lMethod) {
	case loginMethodToken:
		promptI = &survey.Input{Message: "Token: "}
		if err := survey.AskOne(promptI, &token, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	case loginMethodPassword:
		promptI = &survey.Input{Message: "Username: "}
		if err = survey.AskOne(promptI, &user, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		promptPW := &survey.Password{Message: "Password: "}
		if err = survey.AskOne(promptPW, &passwd, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	case loginMethodSSH:
		sshKey, err = prompts.Select("Select SSH-key / -certificate: ", task.ListSSHPubkey(), "[custom filepath]", "")
		if err != nil {
			return err
		}
		fmt.Println(sshKey)

		if strings.Contains(sshKey, "(ssh-agent)") {
			sshAgent = true
			sshKey = ""
		}
		if strings.Contains(sshKey, "principals") {
			sshCertPrincipal = regexp.MustCompile(`.*?principals: (.*?)[,|\s]`).FindStringSubmatch(sshKey)[1]
			if !sshAgent {
				// CLEANUP: should rewrite this with a SSHKey struct with .String() method etc?
				sshKey = regexp.MustCompile(`\((.*?)\)$`).FindStringSubmatch(sshKey)[1]
				sshKey = strings.TrimSuffix(sshKey, "-cert.pub")
				fmt.Println(sshKey)
			}
		} else {
			matches := regexp.MustCompile(`(SHA256:.*?)\s`).FindStringSubmatch(sshKey)
			if len(matches) > 1 {
				sshKeyFingerprint = matches[1]
				if !sshAgent {
					sshKey = regexp.MustCompile(`\((.*?)\)$`).FindStringSubmatch(sshKey)[1]
					sshKey = strings.TrimSuffix(sshKey, ".pub")
				}
			}
		}
	}

	var optSettings bool
	promptYN := &survey.Confirm{
		Message: "Set Optional settings: ",
		Default: false,
	}
	if err = survey.AskOne(promptYN, &optSettings); err != nil {
		return err
	}
	if optSettings {
		// FIXME: drop this prompt entirely, once go's ssh key signing implementation
		//  supports all key types or something??
		//  (at least ecdsa-sha2-nistp521 is unsupported as of 2022-09-03)
		if loginMethod(lMethod) != loginMethodSSH {
			promptI = &survey.Input{Message: "SSH Key Path (leave empty for auto-discovery):"}
			if err := survey.AskOne(promptI, &sshKey); err != nil {
				return err
			}
		}

		promptYN = &survey.Confirm{
			Message: "Allow insecure connections: ",
			Default: false,
		}
		if err = survey.AskOne(promptYN, &insecure); err != nil {
			return err
		}
	}

	return task.CreateLogin(name, token, user, passwd, sshKey, giteaURL, sshCertPrincipal, sshKeyFingerprint, insecure, sshAgent)
}
