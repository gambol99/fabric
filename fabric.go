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

import "github.com/golang/glog"

import (


)

/* the api interface for fabric */
type Fabric interface {
	/* shutdown the service */
	Shutdown()
}

type FabricService struct {
	/* the cluster interface */
	cluster Cluster
	/* the authenticator service */
	auth Authenticator
	/* the docker store client */
	containers ContainerStore
	/* the shutdown channel */
	shutdown ShutdownChannel
}

func NewFabricService() (Fabric, error) {
	fabric := new(FabricService)
	fabric.shutdown = make(ShutdownChannel)
	/* step: create the authenticator service */
	if auth, err := NewAuthenticator(); err != nil {
		return nil, err
	} else {
		fabric.auth = auth
	}
	/* step: create the docker store */
	if client, err := NewContainerStore(); err != nil {
		return nil, err
	} else {
		fabric.containers = client
	}

	/* step: we need to connect to the cluster */
	if cluster, err := NewCluster(); err != nil {
		return nil, err
	} else {
		fabric.cluster = cluster
	}
	return fabric, nil
}

func (r FabricService) Shutdown() {
	glog.Infof("Shutting down the %s service", PROG)
}
