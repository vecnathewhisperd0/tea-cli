// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package task

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
)

// SetupHelper add tea helper to config global
func SetupHelper(login config.Login) (ok bool, err error) {
	// Check that the URL is not blank
	if login.URL == "" {
		return false, fmt.Errorf("Invalid gitea url")
	}

	// get tea binary path
	var binPath string
	if binPath, err = os.Executable(); err != nil {
		return
	}

	// get all helper to URL in git config
	var currentHelpers []byte
	if currentHelpers, err = exec.Command("git", "config", "--global", "--get-all", fmt.Sprintf("credential.%s.helper", login.URL)).Output(); err != nil {
		return false, err
	}

	// Check if ared added tea helper
	for _, line := range strings.Split(strings.ReplaceAll(string(currentHelpers), "\r", ""), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		} else if strings.HasPrefix(line, binPath) && strings.Contains(line[len(binPath):], "login helper") {
			return false, nil
		}
	}

	// Check if tea path have space, if have add quotes
	if strings.Contains(binPath, " ") {
		binPath = fmt.Sprintf("%q", binPath)
	}

	// Add tea helper
	if _, err = exec.Command("git", "config", "--global", "--add", fmt.Sprintf("credential.%s.helper", login.URL), fmt.Sprintf("!%s login helper", binPath)).Output(); err != nil {
		return false, err
	}

	return true, nil
}

// CreateLogin create a login to be stored in config
func CreateLogin(name, token, user, passwd, otp, scopes, sshKey, giteaURL, sshCertPrincipal, sshKeyFingerprint string, insecure, sshAgent, versionCheck, addHelper bool) error {
	// checks ...
	// ... if we have a url
	if len(giteaURL) == 0 {
		return fmt.Errorf("You have to input Gitea server URL")
	}

	// ... if there already exist a login with same name
	if login := config.GetLoginByName(name); login != nil {
		return fmt.Errorf("login name '%s' has already been used", login.Name)
	}
	// ... if we already use this token
	if login := config.GetLoginByToken(token); login != nil {
		return fmt.Errorf("token already been used, delete login '%s' first", login.Name)
	}

	if !sshAgent && sshCertPrincipal == "" && sshKey == "" {
		// .. if we have enough information to authenticate
		if len(token) == 0 && (len(user)+len(passwd)) == 0 {
			return fmt.Errorf("No token set")
		} else if len(user) != 0 && len(passwd) == 0 {
			return fmt.Errorf("No password set")
		} else if len(user) == 0 && len(passwd) != 0 {
			return fmt.Errorf("No user set")
		}
	}

	// Normalize URL
	serverURL, err := utils.NormalizeURL(giteaURL)
	if err != nil {
		return fmt.Errorf("Unable to parse URL: %s", err)
	}

	// check if it's a certificate the principal doesn't matter as the user
	// has explicitly selected this private key
	if _, err := os.Stat(sshKey + "-cert.pub"); err == nil {
		sshCertPrincipal = "yes"
	}

	login := config.Login{
		Name:              name,
		URL:               serverURL.String(),
		Token:             token,
		Insecure:          insecure,
		SSHKey:            sshKey,
		SSHCertPrincipal:  sshCertPrincipal,
		SSHKeyFingerprint: sshKeyFingerprint,
		SSHAgent:          sshAgent,
		Created:           time.Now().Unix(),
		VersionCheck:      versionCheck,
	}

	if len(token) == 0 && sshCertPrincipal == "" && !sshAgent && sshKey == "" {
		if login.Token, err = generateToken(login, user, passwd, otp, scopes); err != nil {
			return err
		}
	}

	client := login.Client()

	// Verify if authentication works and get user info
	u, _, err := client.GetMyUserInfo()
	if err != nil {
		return err
	}
	login.User = u.UserName

	if len(login.Name) == 0 {
		if login.Name, err = GenerateLoginName(giteaURL, login.User); err != nil {
			return err
		}
	}

	// we do not have a method to get SSH config from api,
	// so we just use the host
	login.SSHHost = serverURL.Host

	if len(sshKey) == 0 {
		login.SSHKey, err = findSSHKey(client)
		if err != nil {
			fmt.Printf("Warning: problem while finding a SSH key: %s\n", err)
		}
	}

	if err = config.AddLogin(&login); err != nil {
		return err
	}

	fmt.Printf("Login as %s on %s successful. Added this login as %s\n", login.User, login.URL, login.Name)
	if addHelper {
		if _, err := SetupHelper(login); err != nil {
			return err
		}
	}

	return nil
}

// generateToken creates a new token when given BasicAuth credentials
func generateToken(login config.Login, user, pass, otp, scopes string) (string, error) {
	opts := []gitea.ClientOption{gitea.SetBasicAuth(user, pass)}
	if otp != "" {
		opts = append(opts, gitea.SetOTP(otp))
	}
	client := login.Client(opts...)

	tl, _, err := client.ListAccessTokens(gitea.ListAccessTokensOptions{
		ListOptions: gitea.ListOptions{Page: -1},
	})
	if err != nil {
		return "", err
	}
	host, _ := os.Hostname()
	tokenName := host + "-tea"

	// append timestamp, if a token with this hostname already exists
	for i := range tl {
		if tl[i].Name == tokenName {
			tokenName += time.Now().Format("2006-01-02_15-04-05")
			break
		}
	}

	var tokenScopes []gitea.AccessTokenScope
	for _, scope := range strings.Split(scopes, ",") {
		tokenScopes = append(tokenScopes, gitea.AccessTokenScope(strings.TrimSpace(scope)))
	}

	t, _, err := client.CreateAccessToken(gitea.CreateAccessTokenOption{
		Name:   tokenName,
		Scopes: tokenScopes,
	})
	return t.Token, err
}

// GenerateLoginName generates a name string based on instance URL & adds username if the result is not unique
func GenerateLoginName(url, user string) (string, error) {
	parsedURL, err := utils.NormalizeURL(url)
	if err != nil {
		return "", err
	}
	name := parsedURL.Host

	// append user name if login name already exists
	if len(user) != 0 {
		if login := config.GetLoginByName(name); login != nil {
			return name + "_" + user, nil
		}
	}

	return name, nil
}
