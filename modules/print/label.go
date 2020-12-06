// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"

	"code.gitea.io/tea/cmd/flags"

	"code.gitea.io/sdk/gitea"
	"github.com/muesli/termenv"
)

// LabelsList prints a listing of labels
func LabelsList(labels []*gitea.Label) {
	var values [][]string
	headers := []string{
		"Index",
		"Color",
		"Name",
		"Description",
	}

	if len(labels) == 0 {
		OutputList(flags.GlobalOutputValue, headers, values)
		return
	}

	p := termenv.ColorProfile()

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
	OutputList(flags.GlobalOutputValue, headers, values)
}
