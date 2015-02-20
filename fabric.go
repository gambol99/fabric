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
	"errors"
	"strings"

	swim "github.com/hashicorp/memberlist"
	"github.com/golang/glog"
)

/* the api interface for fabric */
type Fabric interface {
	/* shutdown the service */
	Shutdown()
}

type FabricService struct {
	swim.EventDelegate
	// The cluster configuration for the SWIM protocol
	cluster_config *swim.Config
	// The SWIM client for membership knowledge; note, SWIM is only used to provide a
	// reliable means of cluster membership and notification of membership changes.
	// Messages and state exchanged between members is performed by google protocol
	// buffers
	cluster *swim.Memberlist
	/* the authenticator service */
	auth Authenticator
	/* the docker store client */
	containers ContainerStore
	/* the shutdown channel */
	shutdown ShutdownChannel
	/* the members state */
}

func NewFabricService() (Fabric, error) {
	var err error
	fabric := new(FabricService)
	fabric.shutdown = make(ShutdownChannel)
	/* step: create the authenticator service */
	fabric.auth, err = NewAuthenticator(); Assert(err)
	fabric.containers, err = NewContainerStore(); Assert(err)
	/* step: we need to setup the cluster membership */
	if err := fabric.SetupClusterMembership(); err != nil {
		glog.Errorf("Failed to setup the cluster, error: %s", err)
		return nil, err
	}
	return fabric, nil
}

func (r *FabricService) SetupClusterMembership() error {
	var err error
	glog.Infof("Initializing the cluster membership, boostrap: %t", Options.Bootstrap)
	/* step; we need to generate the configuration */
	r.cluster_config, err = r.ClusterConfiguration(); Assert(err)
	/* step: create the memberlist client */
	r.cluster, err = swim.Create(r.cluster_config); Assert(err)
	/* step: are we bootstrapping?? or join */
	if Options.Bootstrap {
		glog.Infof("No need to join cluster, as we are bootstrapping the cluster")
	} else {
		/* step: extract the members hostname / ipaddress and validate them */
		members := strings.Split(Options.Members, ",")
		for _, member := range members {
			if !IsValidHost(member) {
				glog.Errorf("The member: %s is invalid, please recheck", member)
			}
		}
		glog.Infof("Attempting to join %d cluster members: %s", len(members), Options.Members)
		if successful, err := r.cluster.Join(members); err != nil {
			glog.Errorf("Failed to join cluster members, error: %s", err)
			return err
		} else {
			glog.Infof("Successfully joined with %d %s members", successful, PROG)
		}

	}
	return nil
}

// NotifyJoin is invoked when a node is detected to have joined.
// The Node argument must not be modified.
func (r *FabricService) NotifyJoin(node *swim.Node) {
	glog.V(4).Infof("Member join event, node: %s", node)


}

// NotifyLeave is invoked when a node is detected to have left.
// The Node argument must not be modified.
func (r *FabricService) NotifyLeave(node *swim.Node) {
	glog.V(4).Infof("Member left event, node: %s", node)

}

// NotifyUpdate is invoked when a node is detected to have
// updated, usually involving the meta data. The Node argument
// must not be modified.
func (r *FabricService) NotifyUpdate(node *swim.Node) {
	glog.V(4).Infof("Member update event, node: %s", node)

}

func (r *FabricService) ClusterConfiguration() (*swim.Config, error ) {
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
	/* step: events and delegation */
	config.Events = r
	return config, nil
}

func (r FabricService) Shutdown() {
	glog.Infof("Shutting down the %s service", PROG)
}
