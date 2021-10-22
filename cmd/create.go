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
	"github.com/tabarnhack/git-switch/io/print"
	"github.com/tabarnhack/git-switch/io/prompt/user"
)

var autoAdd bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new git profile",
	Long: `This will create a new git profile and will store it
inside the git profiles DB. Each of the profile must
have a unique name in order to differentiate each
other. The newly created git profile can also be 
automatically added as the current git profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		currUser, err := user.CreateUser(currUser, false)
		if err != nil {
			print.Error("Can't get user information:", err)
			os.Exit(1)
		}

		err = usersDB.Add(currUser)
		if err != nil {
			print.Error("Can't add user to database:", err)
			os.Exit(1)
		}

		err = usersDB.Save()
		if err != nil {
			print.Error("Can't save new user:", err)
			os.Exit(1)
		}

		print.Success("User created")

		if autoAdd {
			switchCmd.Run(cmd, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVar(&currUser.Name, "name", "", "new user's name")
	createCmd.PersistentFlags().StringVar(&currUser.Email, "email", "", "new user's email")
	createCmd.PersistentFlags().BoolVarP(&autoAdd, "auto-add", "a", false, "automatically switch the profile to the one created")
}
