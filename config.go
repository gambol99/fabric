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
	"flag"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var Options struct {
	/* a list of endpoints to connect to */
	Members string `yaml:"members"`
	/* the port the service should be listening on */
	Members_Port int `yaml:"member_port"`
	/* a configuration file containing all of these options */
	ConfigFile string
	/* the ssh key to use */
	Key_File string `yaml:"ssh_ida"`
	/* the ssh port to use */
	Service_Port int `yaml:"port"`
}

func init() {
	flag.StringVar(&Options.ConfigFile, "config", "", "a configuration file container the options")
	flag.StringVar(&Options.Members, "members", "", "a comma seperated list of members to connect")
	flag.IntVar(&Options.Members_Port, "port", 7369, "the port the service should be listening on")
	flag.StringVar(&Options.Key_File, "key-file", "", "the ssh private key to use for the service")
	flag.IntVar(&Options.Service_Port, "ssh-port", 1022, "the port the service should run on")
}

func LoadConfig() error {
	/* step: parse the command line options */
	flag.Parse()
	/* step: check if we're loading a config file */
	if Options.ConfigFile != "" {
		/* step: check the file exists */
		if _, err := os.Stat(Options.ConfigFile); os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("The config file: %s does not exists, please check", Options.ConfigFile))
		}
		/* step: load the content of the file */
		if content, err := ioutil.ReadFile(Options.ConfigFile); err != nil {
			return errors.New(fmt.Sprintf("Failed to read the config file: %s, error: %s", Options.ConfigFile, err))
		} else {
			/* step: attempt to unmarshall the options */
			if err := yaml.Unmarshal([]byte(content), &Options); err != nil {
				return errors.New(fmt.Sprintf("Failed to read the config file: %s, error: %s", Options.ConfigFile, err))
			}
		}
	}
	return nil
}
