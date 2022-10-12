/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

type Product struct {
	Id   int64  `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

const CREATE_PRODUCT = `CREATE TABLE IF NOT EXISTS products (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

const INSERT_PRODUCT = "INSERT INTO products(name) VALUES (?)"

const SELECT_PRODUCTS = "SELECT id, name FROM products;"

const UPDATE_PRODUCT = "UPDATE products SET name = ? WHERE id = ?"

const DELETE_PRODUCT = "DELETE FROM products WHERE id = ?"
