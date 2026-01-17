/*
Copyright (c) 2022-2026, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"fmt"
	"log"

	"github.com/TestAndWin/e2e-coverage/dependency"
	_ "github.com/TestAndWin/e2e-coverage/docs"
	"github.com/TestAndWin/e2e-coverage/router"
)

// @title e2ecoverage
// @version 1.0
// @description API for e2e-coverage

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	fmt.Println("**** e2e test coverage ****")

	// Initialize dependency container
	container := dependency.GetContainer()
	defer func() {
		log.Println("Closing database connections...")
		container.CloseConnections()
	}()

	// Start the router
	router.HandleRequest()
}
