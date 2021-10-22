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
package base

type Entry struct {
	Name  string
	Email string
}

func (e Entry) String() string {
	return e.Name + ": " + e.Email
}

func (e Entry) IsEmpty() bool {
	return e.Name == "" && e.Email == ""
}

func (e Entry) IsIncomplete() bool {
	return e.Name == "" || e.Email == ""
}
