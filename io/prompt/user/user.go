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
package user

import (
	"net/mail"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/tabarnhack/git-switch/base"
	"github.com/tabarnhack/git-switch/io/prompt"
)

func CreateUser(prev base.Entry, isEdit bool) (base.Entry, error) {
	var err error
	if prev.Name == "" || isEdit {
		prev.Name, err = prompt.PromptString("Name", prev.Name, nil)
		if err != nil {
			return base.Entry{}, err
		}
	}

	if prev.Email == "" || isEdit {
		prev.Email, err = prompt.PromptString("Email", prev.Email, func(input string) error {
			_, err := mail.ParseAddress(input)
			return err
		})
		if err != nil {
			return base.Entry{}, err
		}
	}

	return prev, nil
}

func SelectUser(usersDB *base.Base, msg string) (base.Entry, error) {
	entries := usersDB.List()

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "→ {{ .Name | cyan | bold }}: {{ .Email | red }}",
		Inactive: "  {{ .Name | cyan }}: {{ .Email | red }}",
		Selected: "→ {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		entry := entries[index]
		name := strings.Replace(strings.ToLower(entry.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Git profiles",
		Items:     entries,
		Templates: templates,
		Size:      5,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return base.Entry{}, err
	}

	return entries[i], nil
}
