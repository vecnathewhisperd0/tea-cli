// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package interact

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"code.gitea.io/tea/modules/task"
	"golang.org/x/oauth2"

	"github.com/AlecAivazis/survey/v2"
)

// CreateLogin create an login interactive
func CreateLogin() error {
	var (
		name, token, user, passwd, sshKey, giteaURL, sshCertPrincipal, sshKeyFingerprint string
		insecure, sshAgent, versionCheck                                                 bool
	)

	versionCheck = true

	promptI := &survey.Input{Message: "URL of Gitea instance: "}
	if err := survey.AskOne(promptI, &giteaURL, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	giteaURL = strings.TrimSuffix(strings.TrimSpace(giteaURL), "/")
	if len(giteaURL) == 0 {
		fmt.Println("URL is required!")
		return nil
	}

	name, err := task.GenerateLoginName(giteaURL, "")
	if err != nil {
		return err
	}

	oauthToken, err := getToken(giteaURL)
	if err != nil {
		return err
	}

	token = oauthToken.AccessToken

	return task.CreateLogin(name, token, user, passwd, sshKey, giteaURL, sshCertPrincipal, sshKeyFingerprint, oauthToken.RefreshToken, insecure, sshAgent, versionCheck, oauthToken.Expiry)
}

func getToken(giteaURL string) (*oauth2.Token, error) {
	c := oauth2.Config{
		ClientID: "e90ee53c-94e2-48ac-9358-a874fb9e0662",
		Endpoint: oauth2.Endpoint{
			AuthURL:  giteaURL + "/login/oauth/authorize",
			TokenURL: giteaURL + "/login/oauth/access_token",
		},
	}
	state := oauth2.GenerateVerifier()
	queries := make(chan url.Values)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries <- r.URL.Query()
		w.Write([]byte("okay"))
	})
	server := httptest.NewServer(handler)
	c.RedirectURL = server.URL
	defer server.Close()
	verifier := oauth2.GenerateVerifier()
	authCodeURL := c.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	fmt.Fprintf(os.Stderr, "Please complete authentication in your browser...\n%s\n", authCodeURL)
	var open string
	switch runtime.GOOS {
	case "windows":
		open = "start"
	case "darwin":
		open = "open"
	default:
		open = "xdg-open"
	}
	// TODO: wait for server to start before opening browser
	if _, err := exec.LookPath(open); err == nil {
		err = exec.Command(open, authCodeURL).Run()
		if err != nil {
			return nil, err
		}
	}
	query := <-queries
	server.Close()
	if query.Get("state") != state {
		return nil, fmt.Errorf("state mismatch")
	}
	code := query.Get("code")
	return c.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))
}
