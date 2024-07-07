/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

import "time"

type Component struct {
	Name     string    `json:"name"     db:"name"        `
	TestRun  time.Time `json:"test-run" db:"testrun"     `
	Total    int64     `json:"total"`
	Passes   int64     `json:"passes"`
	Pending  int64     `json:"pending"`
	Failures int64     `json:"failures"`
	Skipped  int64     `json:"skipped"`
}
