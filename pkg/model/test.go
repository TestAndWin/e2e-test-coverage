/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

import "time"

type Test struct {
	Id        int64     `db:"id"          json:"id"`
	ProductId int64     `db:"product_id"  json:"product-id"`
	AreaId    int64     `db:"area_id"     json:"area-id"`
	FeatureId int64     `db:"feature_id"  json:"feature-id"`
	Suite     string    `db:"suite"       json:"suite"`
	FileName  string    `db:"file"        json:"file-name"`
	Url       string    `db:"url"         json:"url"`
	Total     int64     `db:"total"       json:"total"`
	Passes    int64     `db:"passes"      json:"passes"`
	Pending   int64     `db:"pending"     json:"pending"`
	Failures  int64     `db:"failures"    json:"failures"`
	Skipped   int64     `db:"skipped"     json:"skipped"`
	Uuid      string    `db:"uuid"         json:"uuid"`
	TestRun   time.Time `db:"testrun"     json:"test-run"`
}

const CREATE_TEST = `CREATE TABLE IF NOT EXISTS tests (
	id INT AUTO_INCREMENT PRIMARY KEY,
	product_id int,
	area_id int,
	feature_id int,
	suite VARCHAR(255),  
	file VARCHAR(255),
	url VARCHAR(500), 
	total int,  
	passes int, 
	pending int, 
	failures int, 
	skipped int, 
	testrun datetime,
	uuid VARCHAR(255),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (feature_id) REFERENCES features(id),
	FOREIGN KEY (feature_id) REFERENCES features(id)
	)`

const INSERT_TEST = "INSERT INTO tests (product_id, area_id, feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun) VALUES (?, ?,?,?,?,?,?,?,?,?,?,?,?)"

const INSERT_TEST_NO_AREA_FEATURE = "INSERT INTO tests (product_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun) VALUES (?,?,?,?,?,?,?,?,?,?,?)"

const DELETE_TEST = "DELETE FROM tests WHERE id = ?"

const SELECT_TESTS_28D = "SELECT id, product_id, area_id, feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun FROM tests WHERE feature_id = ? AND testrun > ? ORDER BY suite, file, testrun DESC;"

const SELECT_TESTS_BY_PRODUCT_28D = "SELECT id, product_id, COALESCE(area_id,0) as area_id, COALESCE(feature_id,0) as feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun FROM tests WHERE product_id = ? AND testrun > ? ORDER BY suite, file, testrun DESC;"
