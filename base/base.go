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

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
	"github.com/tabarnhack/git-switch/config"
	"github.com/tabarnhack/git-switch/io/prompt"
	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("users")

type Base struct {
	filename string

	entries map[string]string
}

func New(conf config.DatabaseConfig) (*Base, error) {
	var path string
	var err error

	if conf.Path != "" {
		path = conf.Path
	} else {
		var found bool
		path, found = searchDB(conf)
		if !found {
			path, err = createDB(conf)
			if err != nil {
				return nil, err
			}
		}
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	entries := make(map[string]string)
	err = db.View(func(tx *bolt.Tx) error {
		tx.Bucket(bucketName).ForEach(func(k, v []byte) error {
			entries[string(k)] = string(v)
			return nil
		})
		return nil
	})

	return &Base{filename: path, entries: entries}, err
}

func searchDB(conf config.DatabaseConfig) (string, bool) {
	for _, v := range conf.SearchPaths {
		path := filepath.Join(v, conf.Filename)
		db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err == nil {
			db.Close()
			return path, true
		}
	}

	return "", false
}

func createDB(conf config.DatabaseConfig) (string, error) {
	var path string
	var created bool

	create, err := prompt.Confirm("No database has been found. Do you want to create one")
	if err != nil {
		return "", err
	}

	if !create {
		return "", fmt.Errorf("Cannot find user database inside the search paths %v", conf.SearchPaths)
	}

	for _, v := range conf.SearchPaths {
		if err := os.MkdirAll(v, 0700); err == nil {
			created = true
			path = filepath.Join(v, conf.Filename)
			break
		}
	}

	if !created {
		return "", fmt.Errorf("Cannot create user database in any of the search paths %v", conf.SearchPaths)
	}

	pterm.Info.Println("User database created at:", path)
	return path, nil
}

func (b *Base) List() []Entry {
	entries := make([]Entry, 0)
	for name, email := range b.entries {
		entries = append(entries, Entry{Name: name, Email: email})
	}

	return entries
}

func (b *Base) Get(name string) (Entry, error) {
	if _, ok := b.entries[name]; !ok {
		return Entry{}, fmt.Errorf("no entry with the name %s exists", name)
	}

	return Entry{Name: name, Email: b.entries[name]}, nil
}

func (b *Base) Add(user Entry) error {
	if user.IsIncomplete() {
		return errors.New("cannot add incomplete user to the database")
	}

	if _, ok := b.entries[user.Name]; ok {
		return fmt.Errorf("an entry with the name %s already exists", user.Name)
	}

	b.entries[user.Name] = user.Email

	return nil
}

func (b *Base) Update(prev, curr Entry) error {
	if _, ok := b.entries[prev.Name]; !ok {
		return fmt.Errorf("no entry with the name %s exists", prev.Name)
	}

	if curr.Name != prev.Name {
		if _, ok := b.entries[curr.Name]; ok {
			return fmt.Errorf("cannot update name to %s as it already exists", curr.Name)
		}
		delete(b.entries, prev.Name)
	}

	b.entries[curr.Name] = curr.Email

	return nil
}

func (b *Base) Delete(name string) error {
	if _, ok := b.entries[name]; !ok {
		return fmt.Errorf("no entry with the name %s exists", name)
	}

	delete(b.entries, name)

	return nil
}

func (b *Base) Print() {
	for _, entry := range b.List() {
		fmt.Printf("name=%s, email=%s\n", entry.Name, entry.Email)
	}
}

func (b *Base) Save() error {
	db, err := bolt.Open(b.filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		// We empty the current base, we assume no modification has been made since `New` call
		err := tx.DeleteBucket(bucketName)
		if err != nil {
			return err
		}
		bucket, err := tx.CreateBucket(bucketName)

		for name, email := range b.entries {
			err = bucket.Put([]byte(name), []byte(email))
			if err != nil {
				break
			}
		}

		return err
	})
}
