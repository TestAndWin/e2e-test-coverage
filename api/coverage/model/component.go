/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

import "time"

type Component struct {
	Name    string    `db:"name"        json:"name"`
	TestRun time.Time `db:"testrun"     json:"test-run"`
}
