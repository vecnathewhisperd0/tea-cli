// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package meta

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"
)

// CmdAutocomplete manages autocompletion
var CmdAutocomplete = cli.Command{
	Name:        "autocomplete",
	Usage:       "Install shell completetion for tea",
	Description: "Install shell completetion for tea",
	ArgsUsage:   "<shell type> (bash, zsh)",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "install",
			Usage: "Persist in shell config instead of printing commands",
		},
	},
	Action: runAutocompleteAdd,
}

func runAutocompleteAdd(ctx *cli.Context) error {
	var file, cmds string
	shell := ctx.Args().First()

	switch shell {
	case "zsh":
		file = "contrib/autocomplete.zsh"
		cmds = "echo 'PROG=tea _CLI_ZSH_AUTOCOMPLETE_HACK=1 source %s' >> ~/.zshrc && source ~/.zshrc"

	case "bash":
		file = "contrib/autocomplete.sh"
		cmds = "echo 'PROG=tea source %s' >> ~/.bashrc && source ~/.bashrc"

	default:
		return fmt.Errorf("Must specify valid shell type")
	}

	destPath, err := xdg.ConfigFile("tea/" + file)
	if err != nil {
		return err
	}
	cmds = fmt.Sprintf(cmds, destPath)

	if err := saveAutoCompleteFile(file, destPath); err != nil {
		return err
	}

	if ctx.Bool("install") {
		fmt.Println("Installing in your shellrc")
		installer := exec.Command(shell, "-c", cmds)
		out, err := installer.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Couldn't run the commands: %s %s", err, out)
		}
	} else {
		fmt.Println("\n# Run the following commands to install autocompletion (or use --install)")
		fmt.Println(cmds)
	}

	return nil
}

func saveAutoCompleteFile(file, destPath string) error {
	url := fmt.Sprintf("https://gitea.com/gitea/tea/raw/branch/master/%s", file)
	fmt.Println("Fetching " + url)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	writer, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, res.Body)
	return err
}
