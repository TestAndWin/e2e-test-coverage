/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package repository

import (
	"context"
	"log"
	"time"

	"github.com/TestAndWin/e2e-coverage/coverage/model"
)

const createProductStmt = `CREATE TABLE IF NOT EXISTS products (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

const insertProductStmt = "INSERT INTO products(name) VALUES (?)"

const updateProductStmt = "UPDATE products SET name = ? WHERE id = ?"

const deleteProductStmt = "DELETE FROM products WHERE id = ?"

func (r CoverageStore) CreateProductsTable() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := r.db.ExecContext(ctx, createProductStmt)
	if err != nil {
		log.Printf("Error %s when creating Products DB table\n", err)
		return err
	}
	return nil
}

func (r CoverageStore) InsertProduct(p model.Product) (int64, error) {
	return r.executeSql(insertProductStmt, p.Name)
}

func (r CoverageStore) UpdateProduct(p model.Product) (int64, error) {
	return r.executeSql(updateProductStmt, p.Name, p.Id)
}

func (r CoverageStore) DeleteProduct(id string) (int64, error) {
	return r.executeSql(deleteProductStmt, id)
}

// Returns all products
func (r CoverageStore) GetAllProducts() ([]model.Product, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name FROM products;")
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var products = []model.Product{}
	for rows.Next() {
		var prd model.Product
		if err := rows.Scan(&prd.Id, &prd.Name); err != nil {
			return products, err
		}
		products = append(products, prd)
	}
	if err := rows.Err(); err != nil {
		return []model.Product{}, err
	}
	return products, nil
}
