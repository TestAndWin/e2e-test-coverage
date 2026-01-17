/*
Copyright (c) 2022-2026, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

import "time"

type Test struct {
	Id             int64     `db:"id"          json:"id"`
	ProductId      int64     `db:"product_id"  json:"product-id"`
	AreaId         int64     `db:"area_id"     json:"area-id"`
	FeatureId      int64     `db:"feature_id"  json:"feature-id"`
	Suite          string    `db:"suite"       json:"suite"`
	FileName       string    `db:"file"        json:"file-name"`
	Component      string    `db:"component"   json:"component"`
	Url            string    `db:"url"         json:"url"`
	Total          int64     `db:"total"       json:"total"`
	Passes         int64     `db:"passes"      json:"passes"`
	Pending        int64     `db:"pending"     json:"pending"`
	Failures       int64     `db:"failures"    json:"failures"`
	Skipped        int64     `db:"skipped"     json:"skipped"`
	Uuid           string    `db:"uuid"        json:"uuid"`
	IsFirst        bool      `db:"is_first"    json:"is-first"`
	TestRun        time.Time `db:"testrun"     json:"test-run"`
	FailedTestRuns int64     `                 json:"failed-test-runs"`
	TotalTestRuns  int64     `                 json:"total-test-runs"`
	FirstTotal     int64     `                 json:"first-total"`
}
