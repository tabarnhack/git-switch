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
package gitconfig

import (
	"gopkg.in/ini.v1"

	"github.com/tabarnhack/git-switch/base"
)

const (
	userSection = "user"

	nameKey  = "name"
	emailKey = "email"
)

type Gitconfig struct {
	filename    string
	cfg         *ini.File
	userSection *ini.Section

	Entry base.Entry
}

func New(filename string, readOnly bool) (*Gitconfig, error) {
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}

	s := cfg.Section(userSection)

	return &Gitconfig{
		filename:    filename,
		cfg:         cfg,
		userSection: s,
		Entry:       base.Entry{Name: s.Key(nameKey).String(), Email: s.Key(emailKey).String()},
	}, nil
}

func (g *Gitconfig) Save() error {
	prev := base.Entry{
		Name:  g.userSection.Key(nameKey).String(),
		Email: g.userSection.Key(emailKey).String(),
	}

	// Not needed to write to file if we have the same profile
	if g.Entry == prev {
		return nil
	}

	g.userSection.Key(nameKey).SetValue(g.Entry.Name)
	g.userSection.Key(emailKey).SetValue(g.Entry.Email)

	return g.cfg.SaveToIndent(g.filename, "\t")
}
