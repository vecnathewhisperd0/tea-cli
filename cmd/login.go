// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli/v2"
)

// CmdLogin represents to login a gitea server.
var CmdLogin = cli.Command{
	Name:        "login",
	Usage:       "Log in to a Gitea server",
	Description: `Log in to a Gitea server`,
	Action:      runLoginList,
	Subcommands: []*cli.Command{
		&cmdLoginList,
		&cmdLoginAdd,
	},
}

// CmdLogin represents to login a gitea server.
var cmdLoginAdd = cli.Command{
	Name:        "add",
	Usage:       "Add a Gitea login",
	Description: `Add a Gitea login`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Login name",
		},
		&cli.StringFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Value:    "https://try.gitea.io",
			EnvVars:  []string{"GITEA_SERVER_URL"},
			Usage:    "Server URL",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t"},
			Value:   "",
			EnvVars: []string{"GITEA_SERVER_TOKEN"},
			Usage:   "Access token. Can be obtained from Settings > Applications",
		},
		&cli.StringFlag{
			Name:    "user",
			Value:   "",
			EnvVars: []string{"GITEA_SERVER_USER"},
			Usage:   "User for basic auth (will create token)",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"pwd"},
			Value:   "",
			EnvVars: []string{"GITEA_SERVER_PASSWORD"},
			Usage:   "Password for basic auth (will create token)",
		},
		&cli.StringFlag{
			Name:    "ssh-key",
			Aliases: []string{"s"},
			Usage:   "Path to a SSH key to use for pull/push operations",
		},
		&cli.BoolFlag{
			Name:    "insecure",
			Aliases: []string{"i"},
			Usage:   "Disable TLS verification",
		},
	},
	Action: runLoginAdd,
}

func runLoginAdd(ctx *cli.Context) error {
	return runLoginAddMain(
		ctx.String("name"),
		ctx.String("token"),
		ctx.String("user"),
		ctx.String("password"),
		ctx.String("ssh-key"),
		ctx.String("url"),
		ctx.Bool("insecure"))
}

func runLoginAddMain(name, token, user, passwd, sshKey, giteaURL string, insecure bool) error {

	if len(giteaURL) == 0 {
		log.Fatal("You have to input Gitea server URL")
	}
	if len(token) == 0 && (len(user)+len(passwd)) == 0 {
		log.Fatal("No token set")
	} else if len(user) != 0 && len(passwd) == 0 {
		log.Fatal("No password set")
	} else if len(user) == 0 && len(passwd) != 0 {
		log.Fatal("No user set")
	}

	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("Unable to load config file " + yamlConfigPath)
	}

	client := gitea.NewClient(giteaURL, token)
	if len(token) == 0 {
		client.SetBasicAuth(user, passwd)
	}
	if insecure {
		cookieJar, _ := cookiejar.New(nil)

		client.SetHTTPClient(&http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		})
	}

	u, err := client.GetMyUserInfo()
	if err != nil {
		log.Fatal(err)
	}

	if len(token) == 0 {
		// create token
		host, _ := os.Hostname()
		tl, err := client.ListAccessTokens(gitea.ListAccessTokensOptions{})
		if err != nil {
			return err
		}
		tokenName := host + "-tea"
		for i := range tl {
			if tl[i].Name == tokenName {
				tokenName += time.Now().Format("2006-01-02_15-04-05")
				break
			}
		}
		t, err := client.CreateAccessToken(gitea.CreateAccessTokenOption{Name: tokenName})
		if err != nil {
			return err
		}
		token = t.Token
	}

	fmt.Println("Login successful! Login name " + u.UserName)

	if len(name) == 0 {
		parsedURL, err := url.Parse(giteaURL)
		if err != nil {
			return err
		}
		name = strings.ReplaceAll(strings.Title(parsedURL.Host), ".", "")
		for _, l := range config.Logins {
			if l.Name == name {
				name += "_" + u.UserName
				break
			}
		}
	}

	err = addLogin(Login{
		Name:     name,
		URL:      giteaURL,
		Token:    token,
		Insecure: insecure,
		SSHKey:   sshKey,
		User:     u.UserName,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = saveConfig(yamlConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// CmdLogin represents to login a gitea server.
var cmdLoginList = cli.Command{
	Name:        "ls",
	Usage:       "List Gitea logins",
	Description: `List Gitea logins`,
	Action:      runLoginList,
	Flags:       []cli.Flag{&OutputFlag},
}

func runLoginList(ctx *cli.Context) error {
	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("Unable to load config file " + yamlConfigPath)
	}

	headers := []string{
		"Name",
		"URL",
		"SSHHost",
		"User",
	}

	var values [][]string

	for _, l := range config.Logins {
		values = append(values, []string{
			l.Name,
			l.URL,
			l.GetSSHHost(),
			l.User,
		})
	}

	Output(outputValue, headers, values)

	return nil
}
