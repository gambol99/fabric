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

import(
	"strings"

	swim "github.com/hashicorp/memberlist"
	"github.com/golang/glog"
	"errors"
	"fmt"
)

const (
	MEMBER_ADDED 	= 0 << iota
	MEMBER_REMOVED
	MEMBER_OFFLINE
)

type Cluster interface {
	/* get a list of the members */
	Members() ([]string, error)
	/* broadcast a message to the cluster */
	Broadcast(message string) error
	/* notify me on a membership change */
	MembershipChange(change int, members string)
	/* notify me on a node */
}

type Swim struct {
	/* the swim client - many thanks to hashicorp, awesome guys! */
	cluster *swim.Memberlist
}

func NewCluster() (Cluster, error) {
	glog.Infof("Initializing the cluster membership, members: %s", Options.Members)
	if config, err := NewClusterConfig(); err != nil {
		return nil, err
	} else {
		service := new(Swim)
		if client, err := swim.Create(config); err != nil {
			return nil, err
		} else {
			service.cluster = client
			/* step: are we bootstrapping?? or join */
			if !Options.Bootstrap {
				glog.Infof("Attempting to join in the cluster, %s", Options.Members)
				if contacted, err := service.cluster.Join(strings.Split(Options.Members, ",")); err != nil {
					glog.Errorf("Failed to join cluster members, error: %s", err)
					return nil, err
				} else {
					glog.Infof("Successfully joined with %d %s members", contacted, PROG)
				}
			}
		}
		return service, nil
	}
}

func (r *Swim) Members() ([]string, error) {
	list := make([]string, 0)
	for _, node := range r.cluster.Members() {
		list = append(list, fmt.Sprintf("%s:%d", node.Addr, node.Port))
	}
	return list, nil
}

func (r *Swim) Broadcast(message string) error {
	return nil
}


func (r *Swim) MembershipChange(change int, members string) {

}

func NewClusterConfig() (*swim.Config, error) {
	var config *swim.Config
	/* step: select the profile */
	switch Options.ClusterProfile {
	case "lan":
		config = swim.DefaultLANConfig()
	case "wan":
		config = swim.DefaultWANConfig()
	case "local":
		config = swim.DefaultLocalConfig()
	default:
		return nil, errors.New("Unsupport cluster profile: "+Options.ClusterProfile)
	}
	/* step: fill in the config */
	config.BindAddr = Options.BindAddr
	config.BindPort = Options.BindPort
	config.AdvertiseAddr = Options.AdvertiseAddr
	config.AdvertisePort = Options.AdvertisePort
	if Options.Secret != "" {
		config.SecretKey = []byte(Options.Secret)
	}
	return config, nil
}


