/*
Copyright 2021.

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
package controllers

import (
	"strings"
	"sync"

	profilesv1alpha1 "github.com/bigkevmcd/profiles-controller/api/v1alpha1"
)

// ProfilesRepository implementations record the seen profiles.
type ProfilesRepository interface {
	Add(profilesv1alpha1.Profile) error
	Search(string) ([]profilesv1alpha1.Profile, error)
}

var _ ProfilesRepository = (*InMemoryProfiles)(nil)

type InMemoryProfiles struct {
	sync.Mutex
	Saved []profilesv1alpha1.Profile
}

func (i *InMemoryProfiles) Add(p profilesv1alpha1.Profile) error {
	i.Lock()
	defer i.Unlock()

	i.Saved = append(i.Saved, p)
	return nil
}

func (i *InMemoryProfiles) Search(q string) ([]profilesv1alpha1.Profile, error) {
	i.Lock()
	defer i.Unlock()

	r := []profilesv1alpha1.Profile{}
	for _, v := range i.Saved {
		if strings.Contains(v.Description, q) {
			r = append(r, v)
		}
	}
	return r, nil
}

func (i *InMemoryProfiles) reset() {
	i.Lock()
	defer i.Unlock()
	i.Saved = []profilesv1alpha1.Profile{}
}
