# <img alt='' src='https://gitea.com/repo-avatars/550-80a3a8c2ab0e2c2d69f296b7f8582485' height="40"/> *T E A*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Release](https://raster.shields.io/badge/dynamic/json.svg?label=release&url=https://gitea.com/api/v1/repos/gitea/tea/releases&query=$[0].tag_name)](https://gitea.com/gitea/tea/releases) [![Build Status](https://drone.gitea.com/api/badges/gitea/tea/status.svg)](https://drone.gitea.com/gitea/tea) [![Join the chat at https://img.shields.io/discord/322538954119184384.svg](https://img.shields.io/discord/322538954119184384.svg)](https://discord.gg/Gitea) [![Go Report Card](https://goreportcard.com/badge/code.gitea.io/tea)](https://goreportcard.com/report/code.gitea.io/tea) [![GoDoc](https://godoc.org/code.gitea.io/tea?status.svg)](https://godoc.org/code.gitea.io/tea)

## The official CLI interface for gitea

This project acts as a command line tool for operating on one or multiple Gitea instances.
It use [code.gitea.io/sdk](https://code.gitea.io/sdk) and interact with the Gitea API

![demo gif](https://dl.gitea.io/screenshots/tea_demo.gif)

## Installation

You can use the prebuilt from [dl.gitea.io](https://dl.gitea.io/tea/)


do it the `go` way:
```sh
go get code.gitea.io/tea
go install code.gitea.io/tea
```


Or if you have `brew` installed, you can install `tea` via:

```sh
brew tap gitea/tap https://gitea.com/gitea/homebrew-gitea
brew install tea
```

## Usage

First of all, you have to create a token on your `personal settings -> application` page of your gitea instance.
Use this token to login with `tea`:

```sh
tea login add --name=try --url=https://try.gitea.io --token=xxxxxx
```

Now you can use the `tea` commands:

```sh
login            Log in to a Gitea server
logout           Log out from a Gitea server
issues           List and create issues
pulls, pull, pr  List open pull requests
releases         Create releases
repos            Operate with repositories
labels           Manage issue labels
times, time      Operate on tracked times of a repositorys issues and pulls
open             Open something of the repository on web browser
```

To fetch issues from different repos, use the `--remote` flag (when inside a gitea repository directory) or `--login` & `--repo` flags.

## Compilation

Make sure you have installed a current go version.
To compile the sources yourself run the following:

```sh
git clone https://gitea.com/gitea/tea.git
cd tea
make
```

## Contributing

Fork -> Patch -> Push -> Pull Request

- `make test` run testsuite
- `make vendor` when adding new dependencies
- ... (for other development tasks, check the `Makefile`)

## Authors

* [Maintainers](https://github.com/orgs/go-gitea/people)
* [Contributors](https://github.com/go-gitea/tea/graphs/contributors)

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the
full license text.
