/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package reporter

import "time"

type TestResult struct {
	Area     string
	Feature  string
	Suite    string
	File     string
	Total    int
	Passes   int
	Pending  int
	Failures int
	Skipped  int
	Uuid     string
	TestRun  time.Time
}
