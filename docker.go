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
	"flag"
	"strings"
	"sync"
	"errors"
	"fmt"
	"bytes"
	"io"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
)

type DockerEvents chan string

const (
	CONTAINER_STARTED = 0 << iota
	CONTAINER_DIED
)

const (
	DOCKER_START   = "start"
	DOCKER_DIE     = "die"
	DOCKER_DESTROY = "destroy"
)

var (
	docker_socket string
)

func init() {
	flag.StringVar(&docker_socket, "docker", "/var/run/docker.sock", "the file / path to the docker socket")
}

type ContainerStore interface {
	/* pull a list of containers */
	List() ([]docker.APIContainers, error)
	/* retrieve a information on a specific docker */
	Get(id string) (*docker.Container, error)
	/* check if a container exists */
	Exists(id string) (bool, error)
	/* attach to the container */
	Attach(containerId string, command string, input io.Reader) error
	/* listen for container creations */
	NotifyOnCreation(channel DockerEvents)
	/* listen out for deaths of containers */
	NotifyOnDeath(channel DockerEvents)
}

type Docker struct {
	sync.RWMutex
	/* the docker client */
	client *docker.Client
	/* a lock for the docker events */
	once_lock sync.Once
	/* the shutdown signal */
	shutdown ShutdownChannel
	/* the channel to send creation events */
	creation_events DockerEvents
	/* the channel to send destruction events */
	destruction_events DockerEvents
}

type DockerContainer struct {
	/* the get of the container */
	id string
	/* the name of the container */
	name string
	/* the image of the container */
	image string
	/* the ports which are being exposed */
}

func NewContainerStore() (ContainerStore, error) {
	glog.Infof("Creating a docker store service, socket: %s", docker_socket)
	store := new(Docker)
	/* step: lets create the docker client */
	if client, err := docker.NewClient("unix://" + docker_socket); err != nil {
		glog.Errorf("Failed to create a docker client, socket: %s, error: %s", docker_socket, err)
		return nil, err
	} else {
		store.client = client
		store.shutdown = make(ShutdownChannel)
		if err := store.client.Ping(); err != nil {
			glog.Errorf("Failed to ping via the docker client, errorr: %s", err)
			return nil, err
		}
		/* step: lets create the docker events */
		if err := store.EventProcessor(); err != nil {
			glog.Errorf("Failed to start the events processor, error: %s", err)
			return nil, err
		}
	}
	return store, nil
}

// Retrieve a list of container current running
func (r *Docker) List() ([]docker.APIContainers, error) {
	if containers, err := r.client.ListContainers(docker.ListContainersOptions{}); err != nil {
		glog.Errorf("Failed to retrieve a list of container from docker, error: %s", err)
		return nil, err
	} else {
	 	return containers, nil
	}
}

func (r *Docker) Get(id string) (*docker.Container, error) {
	if container, err := r.client.InspectContainer(id); err != nil {
		glog.Errorf("Failed to retrieve a container: %s from docker, error: %s", id, err)
		return nil, err
	} else {
		return container, nil
	}
}

func (r *Docker) Exists(id string) (bool, error) {
	if _, err := r.client.InspectContainer(id); err != nil {
		if strings.HasPrefix("No such container", err.Error()) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// The EventProcessor listens out for events from docker and passes them
// upstream to the appropriate channel
func (r *Docker) EventProcessor() error {
	glog.Infof("Starting the Docker Events Processor")

	update_channel := make(chan *docker.APIEvents, 5)
	if err := r.client.AddEventListener(update_channel); err != nil {
		glog.Errorf("Failed to add ourselve as a docker events listener, error: %s", err)
		return err
	}
	/* step: start the events processor */
	go func() {
		glog.Infof("Starting the events processor for docker events")
		for {
			select {
			case event := <- update_channel:
				glog.V(4).Infof("Receivied a docker event, id: %s, status: %s", event.ID[:12], event.Status)
				/* step: are we creating or dying */
				switch event.Status {
				case DOCKER_START:
					if r.creation_events != nil {
						go func() { r.creation_events <- event.ID }()
					}
				case DOCKER_DESTROY:
					if r.destruction_events != nil {
						go func() { r.destruction_events <- event.ID }()
					}
				}
			case <- r.shutdown:
				glog.Infof("Recieved a shutdown signal from above, closing up resources")
				r.client.RemoveEventListener(update_channel)
				break
			}
		}
		glog.Infof("Exitting the events processor loop")
	}()
	return nil
}

// Notify the channel when a container has been created
// Params:
//		channel:	the channel to send the event upon
func (r *Docker) NotifyOnCreation(channel DockerEvents) {
  	r.Lock()
	defer r.Unlock()
	glog.V(6).Infof("Setting the channel for creation events, channel: %v", channel)
	r.creation_events = channel
}

// Notify the channel when a container has died or is killed
// Params:
//		channel:	the channel to send the event upon
func (r *Docker) NotifyOnDeath(channel DockerEvents) {
	r.Lock()
	defer r.Unlock()
	glog.V(6).Infof("Setting the channel for destruction events, channel: %v", channel)
	r.destruction_events = channel
}

func (r *Docker) Attach(id, command string, input io.Reader) error {
	glog.V(4).Infof("Attempting to attach to the container: %s, command: %s", id, command)
	/* step: check the container exists */
	if found, err := r.Exists(id); err != nil {
		return err
	} else if found == false {
		return errors.New(fmt.Sprintf("The container: %s does not exists", id))
	} else {
		exec_options := &docker.CreateExecOptions{
			Container:    id,
			Cmd:          strings.Split(command, " "),
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
		}
		if input != nil {
			exec_options.AttachStdin = true
			exec_options.Tty = true
		}
		if create_exec, err := r.client.CreateExec(*exec_options); err != nil {
			glog.Errorf("Failed to create exec options: %s, error: %s", exec_options, err)
			return err
		} else {
			glog.V(5).Info("Created exec id: %s for container: %s", create_exec.ID, id[:12])
			/* step: we create the StartExecOptions */
			exec := docker.StartExecOptions{
				OutputStream: new(bytes.Buffer),
				ErrorStream:  new(bytes.Buffer),
				InputStream:  input,
				RawTerminal:  false,
			}

			var _ = exec
		}
		return nil
	}
}
