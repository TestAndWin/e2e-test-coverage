/*
Copyright (c) 2022-2026, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

import "time"

type ExplTest struct {
	Id      int64     `db:"id"          json:"id"`
	AreaId  int64     `db:"area_id"     json:"area-id"`
	Summary string    `db:"summary"     json:"summary"`
	Rating  int64     `db:"rating"      json:"rating"`
	TestRun time.Time `db:"testrun"     json:"test-run" time_format:"2006-01-02"`
	Tester  int64     `db:"tester"      json:"tester"`
}
