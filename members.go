/*
Copyright 2014 Rohith All rights reserved.

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

package main

import (
	"sync"
)

type Members interface {
	// Add a member to the list
	Add(id string, member *Member)
	// Removes a member from the list
	Remove(id string)
	// Retrieve a list of the current members
	List() []*Member
	// The size of the members
	Size() int
}

type ClusterMembers struct {
	sync.RWMutex
	// A map of members which memberlist updates
	Members map[string]*Member
}

func (r *ClusterMembers) Add(id string, member *Member) {
	r.Lock()
	defer r.Unlock()
	r.Members[id] = member
}

func (r *ClusterMembers) Remove(id string) {
	r.Lock()
	defer r.Unlock()
	delete(r.Members, id)
}

func (r *ClusterMembers) List() []string {
	r.RLock()
	defer r.RUnlock()
	list := make([]string, 0)
	for id, _ := range r.Members {
		list = append(list, id)
	}
	return list
}

func (r *ClusterMembers) Size() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.Members)
}
