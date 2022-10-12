/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"log"
	"os"

	"github.com/TestAndWin/e2e-coverage/pkg/repository"
)

var repo = initRepository()

// Set-up the db connection and create the db tables if needed
func initRepository() repository.Repository {
	repository, err := repository.OpenDbConnection()
	if err != nil {
		log.Fatal("Error connecting to DB")
		os.Exit(1)
	}
	err = repository.CreateTables()
	if err != nil {
		log.Fatal("Error connecting to DB")
		os.Exit(1)
	}
	return repository
}
