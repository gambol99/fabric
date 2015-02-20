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
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const (
	SSHD_PORT       = 1022
	BIND_PORT 	    = 7946
	BIND_ADDR 	    = "0.0.0.0"
	ADVERTISE_PORT  = BIND_PORT
	ADVERTISE_ADDR  = ""
	CLUSTER_PROFILE = "lan"
)

var Options struct {
	Bootstrap bool
	/* a list of endpoints to connect to */
	Members string
	/* the cluster profile to use */
	ClusterProfile string
	/* the ip address for cluster */
	BindAddr string
	/* the port the service should be listening on */
	BindPort int
	/* the advertised ip address for cluster */
	AdvertiseAddr string
	/* the advertised port to for cluster */
	AdvertisePort int
	/* a configuration file containing all of these options */
	ConfigFile string
	/* the ssh key to use */
	Key_File string
	/* the ssh port to use */
	Service_Port int
	/* the secret key used for cluster membership */
	Secret string
}

func init() {
	flag.BoolVar(&Options.Bootstrap, "bootstrap", false, "use as the bootstrap node")
	flag.StringVar(&Options.ConfigFile, "config", "", "a configuration file container the options")

	/* -- related to the cluster membership */
	flag.StringVar(&Options.ClusterProfile, "profile", CLUSTER_PROFILE, "the cluster profile we should be using")
	flag.StringVar(&Options.Members, "members", "", "a comma seperated list of members to connect")
	flag.IntVar(&Options.BindPort, "bind_port", BIND_PORT, "the port the service should be listening on")
	flag.IntVar(&Options.AdvertisePort, "advertised_port", ADVERTISE_PORT, "the port we should advertise for cluser membership")
	flag.StringVar(&Options.BindAddr, "bind_address", BIND_ADDR, "the address the service should be listening on")
	flag.StringVar(&Options.AdvertiseAddr, "advertised_address", ADVERTISE_ADDR, "the address we should advertise for cluser membership")
	flag.StringVar(&Options.Secret, "secret", "", "the secret used to protect and exchange cluster messages")

	/* -- related to the ssh service */
	flag.StringVar(&Options.Key_File, "keyfile", "", "the ssh private key to use for the service")
	flag.IntVar(&Options.Service_Port, "sshport", SSHD_PORT, "the port the service should run on")
}

func LoadConfig() error {
	/* step: parse the command line options */
	flag.Parse()
	/* step: check if we're loading a config file */
	if Options.ConfigFile != "" {
		/* step: check the file exists */
		if exists, err := FileExists(Options.ConfigFile); err != nil {
			return err
		} else if !exists {
			return errors.New(fmt.Sprintf("The config file: %s does not exists, please check", Options.ConfigFile))
		}

		/* step: load the content of the file */
		if content, err := ioutil.ReadFile(Options.ConfigFile); err != nil {
			return errors.New(fmt.Sprintf("Failed to read the config file: %s, error: %s", Options.ConfigFile, err))
		} else {
			/* step: attempt to un-marshall the options */
			if err := yaml.Unmarshal([]byte(content), &Options); err != nil {
				return errors.New(fmt.Sprintf("Failed to read the config file: %s, error: %s", Options.ConfigFile, err))
			}
		}
	}
	/* step: validate the configuration */
	return ValidateConfig()
}

func ValidateConfig() error {
	/* step: check we have one or more members */
	if Options.Members == "" && !Options.Bootstrap{
		return errors.New("You have not specified any members to join or the bootstrap options")
	}
	/* step: check we have a key file */
	if exists, err := FileExists(Options.Key_File); err != nil {
		return err
	} else if !exists{
		return errors.New(fmt.Sprintf("The key file: '%s' for the service does not exist", Options.Key_File))
	}
	return nil
}
