// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"fmt"
	"net/url"

	git_transport "github.com/go-git/go-git/v5/plumbing/transport"
	gogit_http "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GetAuthForURL returns the appropriate AuthMethod to be used in Push() / Pull()
// operations depending on the protocol, and prompts the user for credentials if
// necessary.
func GetAuthForURL(remoteURL *url.URL, authToken string) (git_transport.AuthMethod, error) {
	switch remoteURL.Scheme {
	case "http", "https":
		// gitea supports push/pull via app token as username.
		return &gogit_http.BasicAuth{Password: "", Username: authToken}, nil
	}
	return nil, fmt.Errorf("don't know how to handle url scheme %v", remoteURL.Scheme)
}
