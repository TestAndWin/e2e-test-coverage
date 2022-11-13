/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
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
	TestRun time.Time `db:"testrun"     json:"test-run" time_format:"2006-01-02" `
}

const CREATE_EXPL_TEST = `CREATE TABLE IF NOT EXISTS expl_tests (
	id INT AUTO_INCREMENT PRIMARY KEY,
	area_id INT, 
	summary TEXT,
	rating INT,
	testrun datetime,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (area_id) REFERENCES areas(id)
	)`

const INSERT_EXPL_TEST = "INSERT INTO expl_tests (area_id, summary, rating, testrun) VALUES (?,?,?,?);"

const DELETE_EXPL_TEST = "DELETE FROM expl_tests WHERE id = ?"

const SELECT_EXPL_TESTS_28D = "SELECT id, area_id, summary, rating, testrun FROM expl_tests WHERE area_id = ? AND testrun > ?;"
