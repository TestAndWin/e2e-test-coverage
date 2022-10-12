/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

type Area struct {
	Id         int64   `db:"id"           json:"id"`
	ProductId  int64   `db:"product_id"   json:"product-id"`
	Name       string  `db:"name"         json:"name"`
	Total      int64   `json:"total"`
	Passes     int64   `json:"passes"`
	Pending    int64   `json:"pending"`
	Failures   int64   `json:"failures"`
	Skipped    int64   `json:"skipped"`
	ExplTests  int64   `json:"expl-tests"`
	ExplRating float64 `json:"expl-rating"`
}

const CREATE_AREA = `CREATE TABLE IF NOT EXISTS areas (
	id INT AUTO_INCREMENT PRIMARY KEY,
	product_id int,
	name VARCHAR(255),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (product_id) REFERENCES products(id)
	)`

const INSERT_AREA = "INSERT INTO areas(product_id, name) VALUES (?, ?)"

const SELECT_AREAS = "SELECT id, product_id, name FROM areas WHERE product_id = ?;"

const UPDATE_AREA = "UPDATE areas SET name = ? WHERE id = ?"

const DELETE_AREA = "DELETE FROM areas WHERE id = ?"
