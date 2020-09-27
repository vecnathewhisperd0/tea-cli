// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"code.gitea.io/tea/modules/intern"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

// CmdLogin represents to login a gitea server.
var CmdLogin = cli.Command{
	Name:        "login",
	Usage:       "Log in to a Gitea server",
	Description: `Log in to a Gitea server`,
	Action:      runLoginAddInteractive,
	Subcommands: []*cli.Command{
		&cmdLoginList,
		&cmdLoginAdd,
		&cmdLoginEdit,
		&cmdLoginSetDefault,
	},
}

// cmdLoginEdit represents to login a gitea server.
var cmdLoginEdit = cli.Command{
	Name:        "edit",
	Usage:       "Edit Gitea logins",
	Description: `Edit Gitea logins`,
	Action:      runLoginEdit,
	Flags:       []cli.Flag{&OutputFlag},
}

func runLoginEdit(ctx *cli.Context) error {
	return open.Start(intern.GetConfigPath())
}

// cmdLoginSetDefault represents to login a gitea server.
var cmdLoginSetDefault = cli.Command{
	Name:        "default",
	Usage:       "Get or Set Default Login",
	Description: `Get or Set Default Login`,
	ArgsUsage:   "<Login>",
	Action:      runLoginSetDefault,
	Flags:       []cli.Flag{&OutputFlag},
}

func runLoginSetDefault(ctx *cli.Context) error {
	if err := intern.LoadConfig(); err != nil {
		return err
	}
	if ctx.Args().Len() == 0 {
		l, err := intern.GetDefaultLogin()
		if err != nil {
			return err
		}
		fmt.Printf("Default Login: %s\n", l.Name)
		return nil
	}
	loginExist := false
	for i := range intern.Config.Logins {
		intern.Config.Logins[i].Default = false
		if intern.Config.Logins[i].Name == ctx.Args().First() {
			intern.Config.Logins[i].Default = true
			loginExist = true
		}
	}

	if !loginExist {
		return fmt.Errorf("login '%s' not found", ctx.Args().First())
	}

	return intern.SaveConfig()
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
	return intern.AddLogin(
		ctx.String("name"),
		ctx.String("token"),
		ctx.String("user"),
		ctx.String("password"),
		ctx.String("ssh-key"),
		ctx.String("url"),
		ctx.Bool("insecure"))
}

func runLoginAddInteractive(ctx *cli.Context) error {
	var stdin, name, token, user, passwd, sshKey, giteaURL string
	var insecure = false

	fmt.Print("URL of Gitea instance: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	giteaURL = strings.TrimSpace(stdin)
	if len(giteaURL) == 0 {
		fmt.Println("URL is required!")
		return nil
	}

	parsedURL, err := url.Parse(giteaURL)
	if err != nil {
		return err
	}
	name = strings.ReplaceAll(strings.Title(parsedURL.Host), ".", "")

	fmt.Print("Name of new Login [" + name + "]: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	if len(strings.TrimSpace(stdin)) != 0 {
		name = strings.TrimSpace(stdin)
	}

	fmt.Print("Do you have a token [Yes/no]: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	if len(stdin) != 0 && strings.ToLower(stdin[:1]) == "n" {
		fmt.Print("Username: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		user = strings.TrimSpace(stdin)

		fmt.Print("Password: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		passwd = strings.TrimSpace(stdin)
	} else {
		fmt.Print("Token: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		token = strings.TrimSpace(stdin)
	}

	fmt.Print("Set Optional settings [yes/No]: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	if len(stdin) != 0 && strings.ToLower(stdin[:1]) == "y" {
		fmt.Print("SSH Key Path: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		sshKey = strings.TrimSpace(stdin)

		fmt.Print("Allow Insecure connections  [yes/No]: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		insecure = len(stdin) != 0 && strings.ToLower(stdin[:1]) == "y"
	}

	return intern.AddLogin(name, token, user, passwd, sshKey, giteaURL, insecure)
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
	err := intern.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Name",
		"URL",
		"SSHHost",
		"User",
		"Default",
	}

	var values [][]string

	for _, l := range intern.Config.Logins {
		values = append(values, []string{
			l.Name,
			l.URL,
			l.GetSSHHost(),
			l.User,
			fmt.Sprint(l.Default),
		})
	}

	Output(globalOutputValue, headers, values)

	return nil
}
