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
package print

import (
	"github.com/pterm/pterm"
)

type TableData [][]string

func Println(v ...interface{}) {
	pterm.FgCyan.Println(v...)
}

func Info(v ...interface{}) {
	pterm.Info.Println(v...)
}

func Success(v ...interface{}) {
	pterm.Success.Println(v...)
}

func Error(v ...interface{}) {
	pterm.Error.Println(v...)
}

func Table(d TableData) {
	pterm.DefaultTable.WithHasHeader().WithData(d).Render()
}

func Section(v ...interface{}) {
	pterm.DefaultSection.Println(v...)
}
