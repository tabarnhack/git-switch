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
	"github.com/tabarnhack/git-switch/gitconfig"
	"github.com/tabarnhack/git-switch/io/print"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Dump git profiles DB",
	Long: `List every git profile stored inside the DB
and the current git profile set.`,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := gitconfig.New(gitconfigFile, true)
		if err != nil {
			fmt.Println("Can't load gitconfig file:", err)
			os.Exit(1)
		}

		print.Section("Active git profile")
		print.Println(g.Entry)

		print.Section("Git profiles list")
		entries := usersDB.List()
		data := print.TableData{[]string{"Name", "Email"}}
		for _, entry := range entries {
			data = append(data, []string{entry.Name, entry.Email})
		}

		print.Table(data)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
