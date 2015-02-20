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

type Container struct {
	/* the id for the container */
	ID string
	/* the name of the container */
	Name string
	/* the image it is running from */
	Image string
	/* the ports which are exposed */
	Ports map[int]int
	/* additional meta data associated with the container */
	Meta map[string]string
	/* is the container running */
	Running bool
}

func NewContainer() *Container {
	return &Container{
		ID:      "",
		Name:    "",
		Image:   "",
		Ports:   make(map[int]int, 0),
		Meta:    make(map[string]string, 0),
		Running: true,
	}
}
