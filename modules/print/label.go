// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"

	"code.gitea.io/sdk/gitea"
	"github.com/muesli/termenv"
)

// LabelsList prints a listing of labels
func LabelsList(labels []*gitea.Label, output string) {
	t := tableWithHeader(
		"Index",
		"Color",
		"Name",
		"Description",
	)

	for _, label := range labels {
		t.addRow(
			strconv.FormatInt(label.ID, 10),
			formatLabel(label, !isMachineReadable(output), label.Color),
			label.Name,
			label.Description,
		)
	}
	t.print(output)
}

func formatLabel(label *gitea.Label, allowColor bool, text string) string {
	colorProfile := termenv.Ascii
	if allowColor {
		colorProfile = termenv.EnvColorProfile()
	}
	if len(text) == 0 {
		text = label.Name
	}
	styled := termenv.String(text)
	styled = styled.Foreground(colorProfile.Color("#" + label.Color))
	return fmt.Sprint(styled)
}
