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
	"bufio"
	"bytes"
	"flag"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
)

func main() {
	var container string
	flag.StringVar(&container, "c", "test1", "container name")

	flag.Parse()

	createOpts := docker.CreateExecOptions{
		Container:    container,
		Cmd:          []string{"ifconfig"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true}

	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	execObj, err := client.CreateExec(createOpts)
	if err != nil {
		fmt.Println("failed to run in container - Exec setup failed - %v", err)
		os.Exit(1)
	}

	startOpts := docker.StartExecOptions{
		Detach:       false,
		Tty:          true,
		OutputStream: new(bytes.Buffer),
		ErrorStream:  new(bytes.Buffer),
		InputStream:  new(bytes.Buffer),
		RawTerminal:  true,
	}
	errChan := make(chan error, 1)
	go func() {
		fmt.Println("starting the exe")
		errChan <- client.StartExec(execObj.ID, startOpts)
	}()
	// read from stdin
	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				// You may check here if err == io.EOF
				break
			}
			startOpts.InputStream.Read([]byte(line))
		}
	}()

	err = <-errChan
	fmt.Printf("finished the exe : %s\n", startOpts.OutputStream)
	if err != nil {
		fmt.Println("failed to run in container - Exec start failed - %v", err)
		os.Exit(1)
	}
}
