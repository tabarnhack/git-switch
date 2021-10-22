/*
Copyright © 2021 Tabarnhack <tabarnhack@outlook.fr>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tabarnhack/git-switch/base"
	"github.com/tabarnhack/git-switch/config"
	"github.com/tabarnhack/git-switch/io/print"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	SYSTEM_GITCONFIG = "/etc/gitconfig"
	GLOBAL_GITCONFIG = "/.gitconfig"
	LOCAL_GITCONFIG  = "/.git/config"
)

var (
	cfgFile       string
	profilesBase  string
	gitconfigFile string

	systemGitconfig bool
	globalGitConfig bool
	localGitconfig  bool

	usersDB  *base.Base
	currUser base.Entry
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-switch",
	Short: "A brief description of your application",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		print.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/git-switch/config.yml)")

	rootCmd.PersistentFlags().StringVarP(&profilesBase, "db", "d", "", "git profiles database")

	rootCmd.PersistentFlags().StringVar(&gitconfigFile, "gitconfig", "", "gitconfig file to use")
	rootCmd.PersistentFlags().BoolVarP(&systemGitconfig, "system", "s", false, "modify gitconfig at system level (eg. /etc/git/gitconfig)")
	rootCmd.PersistentFlags().BoolVarP(&globalGitConfig, "global", "g", false, "modify gitconfig at global level (eg. $HOME/.gitconfig)")
	rootCmd.PersistentFlags().BoolVarP(&localGitconfig, "local", "l", false, "modify gitconfig at local level (eg. $PWD/.git/config)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	conf, configPaths, err := config.New()
	if err != nil {
		print.Error("Can't initialize config:", err)
		os.Exit(1)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		for _, path := range configPaths {
			viper.AddConfigPath(path)
		}
		viper.SetConfigName(config.ConfigName)
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetConfigType("yml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		print.Error("Can't read config:", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		print.Error("Unable to decode into struct, %v", err)
		os.Exit(1)
	}

	print.Info("Using config file:", viper.ConfigFileUsed())
	fmt.Printf("%#v\n", conf)

	if systemGitconfig {
		if gitconfigFile != "" {
			print.Error("Can't specify multiple gitconfig files")
			os.Exit(1)
		}
		gitconfigFile = SYSTEM_GITCONFIG
	}

	if globalGitConfig {
		if gitconfigFile != "" {
			print.Error("Can't specify multiple gitconfig files")
			os.Exit(1)
		}
		home, err := homedir.Dir()
		if err != nil {
			print.Error("Can't get home directory:", err)
			os.Exit(1)
		}
		gitconfigFile = home + GLOBAL_GITCONFIG
	}

	if localGitconfig {
		if gitconfigFile != "" {
			print.Error("Can't specify multiple gitconfig files")
			os.Exit(1)
		}
		pwd, err := os.Getwd()
		if err != nil {
			print.Error("Can't get working directory:", err)
			os.Exit(1)
		}
		gitconfigFile = pwd + LOCAL_GITCONFIG
	}

	// if gitconfigFile is empty, we load a default value
	if gitconfigFile == "" {
		gitconfigFile = conf.DefaultGitconfig
		print.Info("No gitconfig file provided. Using default:", gitconfigFile)
	}

	if profilesBase != "" {
		conf.Database.Path = profilesBase
	}

	usersDB, err = base.New(conf.Database)
	if err != nil {
		print.Error("Can't load user database:", err)
		os.Exit(1)
	}
}
