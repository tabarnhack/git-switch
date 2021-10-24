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
	"os"

	"github.com/spf13/cobra"
	"github.com/tabarnhack/git-switch/gitconfig"
	"github.com/tabarnhack/git-switch/io/print"
	"github.com/tabarnhack/git-switch/io/prompt"
	"github.com/tabarnhack/git-switch/io/prompt/user"
)

var (
	saveExisting bool
	forceSwitch  bool
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch the git profile used in gitconfig",
	Long: `Change the git profile of the corresponding gitconfig
file with the one selected from the DB. The existing
git profile can be saved before being overwritten.`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := gitconfig.New(gitconfigFile, true)
		if err != nil {
			print.Error("Can't load gitconfig file:", err)
			os.Exit(1)
		}

		if g.Entry.IsEmpty() {
			print.Info("Currently, no git profile is set inside this file")
		} else {
			print.Info("The current profile for this gitconfig file is", g.Entry)
			if !saveExisting && !forceSwitch {
				saveExisting, err = prompt.Confirm("Do you want to save the current git profile")
				if err != nil {
					print.Error("Cannot get user confirmation:", err)
					os.Exit(1)
				}
			}

			if saveExisting {
				err = usersDB.Add(g.Entry)
				if err != nil {
					print.Error("Can't add user to database:", err)
					os.Exit(1)
				}
			}
		}

		if currUser.Name == "" {
			currUser, err = user.SelectUser(usersDB, "Switch user")
		} else {
			currUser, err = usersDB.Get(currUser.Name)
		}

		if err != nil {
			print.Error("Can't select user:", err)
			os.Exit(1)
		}

		g.Entry = currUser

		err = g.Save()
		if err != nil {
			print.Error("Can't save edited gitconfig file:", err)
			os.Exit(1)
		}

		print.Success("Selected user:", currUser)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	switchCmd.PersistentFlags().StringVarP(&currUser.Name, "name", "n", "", "name of the user to switch to")
	switchCmd.PersistentFlags().BoolVarP(&saveExisting, "save", "w", false, "save the existing git profile before switching")
	switchCmd.PersistentFlags().BoolVarP(&forceSwitch, "force", "f", false, "force git profile overwrite")
}
