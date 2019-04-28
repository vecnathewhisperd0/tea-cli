// Copyright © 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool

// Version holds the current Gitea version
var Version = "0.1.0-dev"

// Tags holds the build tags used
var Tags = ""

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tea",
	Short: "Command line tool to interact with Gitea",
	Long: `Command line tool to interact with Gitea

Tea is an application to interact with the Gitea API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.tea/tea.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "config file (default is $HOME/.tea/tea.yaml)")
	if len(Tags) > 0 {
		Version += " built with: " + strings.Replace(Tags, " ", ", ", -1)
	}
	rootCmd.Version = Version
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tea" (without extension).
		viper.AddConfigPath(home + "/.tea")
		viper.SetConfigName("tea")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

func getOwnerRepo() (string, string) {
	var err error
	repoPath := rootCmd.PersistentFlags().Lookup("repo").Value.String()
	if repoPath == "" {
		_, repoPath, err = curGitRepoPath()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	owner, repo := splitRepo(repoPath)
	return owner, repo
}
