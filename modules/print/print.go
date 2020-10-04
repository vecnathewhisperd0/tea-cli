// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"time"

	"github.com/muesli/termenv"
)

var (
	location *time.Location
)

func init() {
	// get local timezone for time formatting. we ignore the error
	// and handle that case in FormatTime().
	location, _ = time.LoadLocation("Local")
}

func getGlamourTheme() string {
	if termenv.HasDarkBackground() {
		return "dark"
	}
	return "light"
}

// formatSize get kb in int and return string
func formatSize(kb int64) string {
	if kb < 1024 {
		return fmt.Sprintf("%d Kb", kb)
	}
	mb := kb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%d Mb", mb)
	}
	gb := mb / 1024
	if gb < 1024 {
		return fmt.Sprintf("%d Gb", gb)
	}
	return fmt.Sprintf("%d Tb", gb/1024)
}

// FormatTime give a date-time in local timezone if available
func FormatTime(t time.Time) string {
	if location == nil {
		return t.Format("2006-01-02 15:04 UTC")
	}
	return t.In(location).Format("2006-01-02 15:04")
}
