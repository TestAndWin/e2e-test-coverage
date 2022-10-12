/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"fmt"

	_ "github.com/TestAndWin/e2e-coverage/docs"

	"github.com/TestAndWin/e2e-coverage/pkg/router"
)

// @title e2ecoverage
// @version 1.0
// @description API for e2e-coverage

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	fmt.Println("**** e2e test coverage ****")

	router.HandleRequest()
}
