/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

type Feature struct {
	Id            int64  `json:"id"              db:"id"`
	AreaId        int64  `json:"area-id"         db:"area_id"`
	Name          string `json:"name"            db:"name"`
	Documentation string `json:"documentation"   db:"documentation"`
	Url           string `json:"url"             db:"url"`
	BusinessValue string `json:"business-value"  db:"business_value"`
	Total         int64  `json:"total"`
	FirstTotal    int64  `json:"first-total"`
	Passes        int64  `json:"passes"`
	Pending       int64  `json:"pending"`
	Failures      int64  `json:"failures"`
	Skipped       int64  `json:"skipped"`
	Tests         []Test `json:"tests"`
}
