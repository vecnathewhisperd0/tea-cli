// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"code.gitea.io/tea/modules/intern"

	"code.gitea.io/sdk/gitea"
	"github.com/muesli/termenv"
	"github.com/urfave/cli/v2"
)

// CmdLabels represents to operate repositories' labels.
var CmdLabels = cli.Command{
	Name:        "labels",
	Usage:       "Manage issue labels",
	Description: `Manage issue labels`,
	Action:      runLabels,
	Subcommands: []*cli.Command{
		&CmdLabelCreate,
		&CmdLabelUpdate,
		&CmdLabelDelete,
	},
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "save",
			Aliases: []string{"s"},
			Usage:   "Save all the labels as a file",
		},
		&PaginationPageFlag,
		&PaginationLimitFlag,
	}, AllDefaultFlags...),
}

func runLabels(ctx *cli.Context) error {
	login, owner, repo := intern.InitCommand(globalRepoValue, globalLoginValue, globalRemoteValue)

	headers := []string{
		"Index",
		"Color",
		"Name",
		"Description",
	}

	var values [][]string

	labels, _, err := login.Client().ListRepoLabels(owner, repo, gitea.ListLabelsOptions{ListOptions: getListOptions(ctx)})
	if err != nil {
		log.Fatal(err)
	}

	if len(labels) == 0 {
		Output(globalOutputValue, headers, values)
		return nil
	}

	p := termenv.ColorProfile()

	fPath := ctx.String("save")
	if len(fPath) > 0 {
		f, err := os.Create(fPath)
		if err != nil {
			return err
		}
		defer f.Close()

		for _, label := range labels {
			fmt.Fprintf(f, "#%s %s\n", label.Color, label.Name)
		}
	} else {
		for _, label := range labels {
			color := termenv.String(label.Color)

			values = append(
				values,
				[]string{
					strconv.FormatInt(label.ID, 10),
					fmt.Sprint(color.Background(p.Color("#" + label.Color))),
					label.Name,
					label.Description,
				},
			)
		}
		Output(globalOutputValue, headers, values)
	}

	return nil
}

func splitLabelLine(line string) (string, string, string) {
	fields := strings.SplitN(line, ";", 2)
	var color, name, description string
	if len(fields) < 1 {
		return "", "", ""
	} else if len(fields) >= 2 {
		description = strings.TrimSpace(fields[1])
	}
	fields = strings.Fields(fields[0])
	if len(fields) <= 0 {
		return "", "", ""
	}
	color = fields[0]
	if len(fields) == 2 {
		name = fields[1]
	} else if len(fields) > 2 {
		name = strings.Join(fields[1:], " ")
	}
	return color, name, description
}
