package git

import (
	"gopkg.in/src-d/go-git.v4"
	git_config "gopkg.in/src-d/go-git.v4/config"
)

// GetOrCreateRemote tries to match a Remote of the repo via the given URL.
// If no match is found, a new Remote with `newRemoteName` is created.
// Matching is based on the normalized URL, accepting different protocols.
func (r TeaRepo) GetOrCreateRemote(remoteURL, newRemoteName string) (*git.Remote, error) {
	repoURL, err := ParseURL(remoteURL)
	if err != nil {
		return nil, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}
	var localRemote *git.Remote = nil
	for _, r := range remotes {
		for _, u := range r.Config().URLs {
			remoteURL, _ := ParseURL(u)
			if remoteURL.Host == repoURL.Host && remoteURL.Path == repoURL.Path {
				localRemote = r
				break
			}
		}
		if localRemote != nil {
			break
		}
	}

	// if no match found, create a new remote
	if localRemote == nil {
		localRemote, err = r.CreateRemote(&git_config.RemoteConfig{
			Name: newRemoteName,
			URLs: []string{remoteURL},
		})
		if err != nil {
			return nil, err
		}
	}

	return localRemote, nil
}
