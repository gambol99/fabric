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
)

var Options struct {
	/* a list of endpoints to connect to */
	Members string
	/* the port the service should be listening on */
	Port int
	/* a configuration file containing all of these options */
	ConfigFile string
	}

func init() {
	flag.StringVar(&Options.ConfigFile, "config", "", "a configuration file container the options")
	flag.StringVar(&Options.Members, "members", "", "a comma seperated list of members to connect")
	flag.IntVar(&Options.Port, "port", 1022, "the port the service should be listening on")
}

func LoadConfig() error {
	/* parse the command line options */
	flag.Parse()
	if Options.ConfigFile != "" {

	} else {



	}
	return nil
}
