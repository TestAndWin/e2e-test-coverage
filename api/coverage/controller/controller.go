/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"github.com/TestAndWin/e2e-coverage/coverage/repository"
	"github.com/TestAndWin/e2e-coverage/dependency"
)

// getRepository returns the coverage repository from the dependency container
func getRepository() (*repository.CoverageStore, error) {
	container := dependency.GetContainer()
	repo, err := container.GetCoverageStore()
	if err != nil {
		return nil, err
	}
	return repo, nil
}
