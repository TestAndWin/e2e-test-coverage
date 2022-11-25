/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
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

const CREATE_FEATURE = `CREATE TABLE IF NOT EXISTS features(
	id INT AUTO_INCREMENT PRIMARY KEY,
	area_id INT, 
	name VARCHAR(255),
	documentation VARCHAR(255),
	url VARCHAR(255),
	business_value VARCHAR(20),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (area_id) REFERENCES areas(id)
	)`

const INSERT_FEATURE = "INSERT INTO features(area_id, name, documentation, url, business_value) VALUES (?,?,?,?,?)"

const SELECT_FEATURES = "SELECT id, area_id, name, documentation, url, business_value FROM features WHERE area_id = ?;"

const UPDATE_FEATURE = "UPDATE features SET name = ?, documentation = ?, url = ?, business_value = ? WHERE id = ?"

const DELETE_FEATURE = "DELETE FROM features WHERE id = ?"
