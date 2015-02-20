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

//ref https://github.com/fsouza/go-dockerclient

func main() {
	var container string
	flag.StringVar(&container, "c", "test1", "container name")

	flag.Parse()

	createOpts := docker.CreateExecOptions{
		Container:    container,
		Cmd:          []string{"ip", "addr"},
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false}

	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	execObj, err := client.CreateExec(createOpts)
	if err != nil {
		fmt.Println("failed to run in container - Exec setup failed - %v", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	wrBuf := bufio.NewWriter(&buf)
	startOpts := docker.StartExecOptions{
		Detach:       false,
		Tty:          false,
		OutputStream: wrBuf,
		ErrorStream:  wrBuf,
		RawTerminal:  false,
	}
	errChan := make(chan error, 1)
	go func() {
		errChan <- client.StartExec(execObj.ID, startOpts)
	}()
	err = <-errChan
	if err != nil {
		fmt.Println("failed to run in container - Exec start failed - %v", err)
		os.Exit(1)
	}
	wrBuf.Flush()
	//data := buf.Bytes()
	fmt.Println("result: ", buf.String())
}
