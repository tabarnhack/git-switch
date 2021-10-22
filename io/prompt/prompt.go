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
package prompt

import (
	"errors"
	"strings"

	"github.com/manifoldco/promptui"
)

func PromptString(msg, defaultVal string, moreValidate func(string) error) (string, error) {
	validate := func(input string) error {
		if len(strings.TrimSpace(input)) == 0 {
			return errors.New("Empty string")
		}

		if moreValidate != nil {
			err := moreValidate(input)
			if err != nil {
				return err
			}
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:     msg,
		Default:   defaultVal,
		AllowEdit: true,
		Validate:  validate,
	}

	result, err := prompt.Run()
	return strings.TrimSpace(result), err
}

func Confirm(msg string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil && err.Error() != "" {
		return false, err
	}

	return result == "y", nil
}
