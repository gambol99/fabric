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

	"github.com/golang/glog"
)

var AuthenticatorOptions struct {
	/* the authenticator to user */
	authenticator string
	}

func init() {
	flag.StringVar(&AuthenticatorOptions.authenticator, "auth", "plain", "the authentication method to use")
}

type Authenticator interface {
	/* check the validity of a user */
	AuthenticateLogin(user, token string) (string, error)
	/* check the user is allow to perform this operation */
	Authenticate(userID string, actionID int) (bool, error)
}

type TestAuthentication struct {}

func (r *TestAuthentication) AuthenticateLogin(user, token string) (string, error) {
	return "11111111", nil
}

func (r *TestAuthentication) Authenticate(userID string, actionID int) (bool, error) {
	return true, nil
}

func NewAuthenticator() (Authenticator, error) {
	glog.V(4).Infof("Initializating the Authenticator Service, type: %s", AuthenticatorOptions.authenticator)
	switch AuthenticatorOptions.authenticator {
	case "plain":
		return &TestAuthentication{}, nil
	}
	return nil, errors.New("The authenticator not found or supported")
}




