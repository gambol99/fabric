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
	"os"

	"github.com/golang/glog"
	"syscall"
	"os/signal"
)

func main() {
	/* step: parse the command line options and load configuration */
	if err := LoadConfig(); err != nil {
		glog.Errorf("Invalid configration: %s", err)
		os.Exit(1)
	}
	glog.Infof("Initialize the %s (%s) service", PROG, VERSION)

	/* step: bind our terminal service and wait for incoming requests */
	if fabric, err := NewFabricService(); err != nil {
		glog.Errorf("Failed to initialize the fabric service, error: %s", err)
		os.Exit(1)
	} else {
		glog.Infof("Waiting for signal to quit")
		/* step: setup the channel for shutdown signals */
		signalChannel := make(chan os.Signal)
		/* step: register the signals */
		signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		/* step: wait on the signal */
		<-signalChannel
		/* shutdown the service */
		fabric.Shutdown()
	}
}
