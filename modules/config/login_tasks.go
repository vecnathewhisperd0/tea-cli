// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"code.gitea.io/tea/modules/utils"
)

// GetDefaultLogin return the default login
func GetDefaultLogin() (*Login, error) {
	if len(Config.Logins) == 0 {
		return nil, errors.New("No available login")
	}
	for _, l := range Config.Logins {
		if l.Default {
			return &l, nil
		}
	}

	return &Config.Logins[0], nil
}

// GetLoginByName get login by name
func GetLoginByName(name string) *Login {
	for _, l := range Config.Logins {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

// AddLogin add login to config ( global var & file)
func AddLogin(name, token, user, passwd, sshKey, giteaURL string, insecure bool) error {
	// checks ...
	// ... if we have a url
	if len(giteaURL) == 0 {
		log.Fatal("You have to input Gitea server URL")
	}

	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range Config.Logins {
		// ... if there already exist a login with same name
		if strings.ToLower(l.Name) == strings.ToLower(name) {
			return fmt.Errorf("login name '%s' has already been used", l.Name)
		}
		// ... if we already use this token
		if l.Token == token {
			return fmt.Errorf("token already been used, delete login '%s' first", l.Name)
		}
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

	login := Login{
		Name:     name,
		URL:      serverURL.String(),
		Token:    token,
		Insecure: insecure,
		SSHKey:   sshKey,
		Created:  time.Now().Unix(),
	}

	if len(token) == 0 {
		login.Token, err = login.GenerateToken(user, passwd)
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

	// save login to global var
	Config.Logins = append(Config.Logins, login)

	// save login to config file
	err = SaveConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Login as %s on %s successful. Added this login as %s\n", login.User, login.URL, login.Name)

	return nil
}

// DeleteLogin delete a login by name
func DeleteLogin(name string) error {
	var idx = -1
	for i, l := range Config.Logins {
		if l.Name == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("can not delete login '%s', does not exist", name)
	}

	Config.Logins = append(Config.Logins[:idx], Config.Logins[idx+1:]...)

	return SaveConfig()
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
		for _, l := range Config.Logins {
			if l.Name == name {
				name += "_" + user
				break
			}
		}
	}

	return name, nil
}
