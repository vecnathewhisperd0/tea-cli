// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"
	"log"
	"os"
	"time"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
)

// CreateLogin create a login to be stored in config
func CreateLogin(name, token, user, passwd, sshKey, giteaURL string, insecure bool) error {
	// checks ...
	// ... if we have a url
	if len(giteaURL) == 0 {
		log.Fatal("You have to input Gitea server URL")
	}

	// ... if there already exist a login with same name
	if login := config.GetLoginByName(name); login != nil {
		return fmt.Errorf("login name '%s' has already been used", login.Name)
	}
	// ... if we already use this token
	if login := config.GetLoginByToken(token); login != nil {
		return fmt.Errorf("token already been used, delete login '%s' first", login.Name)
	}

	// .. if we have enough information to authenticate
	if len(token) == 0 && (len(user)+len(passwd)) == 0 {
		log.Fatal("No token set")
	} else if len(user) != 0 && len(passwd) == 0 {
		log.Fatal("No password set")
	} else if len(user) == 0 && len(passwd) != 0 {
		log.Fatal("No user set")
	}

	// Normalize URL
	serverURL, err := utils.NormalizeURL(giteaURL)
	if err != nil {
		log.Fatal("Unable to parse URL", err)
	}

	login := config.Login{
		Name:     name,
		URL:      serverURL.String(),
		Token:    token,
		Insecure: insecure,
		SSHKey:   sshKey,
		Created:  time.Now().Unix(),
	}

	if len(token) == 0 {
		login.Token, err = generateToken(login.Client(), user, passwd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Verify if authentication works and get user info
	u, _, err := login.Client().GetMyUserInfo()
	if err != nil {
		log.Fatal(err)
	}
	login.User = u.UserName

	if len(login.Name) == 0 {
		login.Name, err = GenerateLoginName(giteaURL, login.User)
		if err != nil {
			log.Fatal(err)
		}
	}

	// we do not have a method to get SSH config from api,
	// so we just use the hostname
	login.SSHHost = serverURL.Hostname()

	err = config.AddLogin(&login)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Login as %s on %s successful. Added this login as %s\n", login.User, login.URL, login.Name)

	return nil
}

// generateToken creates a new token when given BasicAuth credentials
func generateToken(client *gitea.Client, user, pass string) (string, error) {
	gitea.SetBasicAuth(user, pass)(client)

	host, _ := os.Hostname()
	tl, _, err := client.ListAccessTokens(gitea.ListAccessTokensOptions{})
	if err != nil {
		return "", err
	}
	tokenName := host + "-tea"

	for i := range tl {
		if tl[i].Name == tokenName {
			tokenName += time.Now().Format("2006-01-02_15-04-05")
			break
		}
	}

	t, _, err := client.CreateAccessToken(gitea.CreateAccessTokenOption{Name: tokenName})
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
