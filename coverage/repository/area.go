/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
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

const createAreaStmt = `CREATE TABLE IF NOT EXISTS areas (
	id INT AUTO_INCREMENT PRIMARY KEY,
	product_id int,
	name VARCHAR(255),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (product_id) REFERENCES products(id)
	)`

const insertAreaStmt = "INSERT INTO areas(product_id, name) VALUES (?, ?)"

const updateAreaStmt = "UPDATE areas SET name = ? WHERE id = ?"

const deleteAreaStmt = "DELETE FROM areas WHERE id = ?"

func (r Repository) CreateAreasTable() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := r.db.ExecContext(ctx, createAreaStmt)
	if err != nil {
		log.Printf("Error %s when creating Areas DB table\n", err)
		return err
	}
	return nil
}

func (r Repository) InsertArea(a model.Area) (int64, error) {
	return r.executeSql(insertAreaStmt, a.ProductId, a.Name)
}

func (r Repository) UpdateArea(a model.Area) (int64, error) {
	return r.executeSql(updateAreaStmt, a.Name, a.Id)
}

func (r Repository) DeleteArea(id string) (int64, error) {
	return r.executeSql(deleteAreaStmt, id)
}

// Get all areas for the specified product id
func (r Repository) GetAllProductAreas(pid string) ([]model.Area, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, product_id, name FROM areas WHERE product_id = ?;")
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, pid)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var areas = []model.Area{}
	for rows.Next() {
		var prd model.Area
		if err := rows.Scan(&prd.Id, &prd.ProductId, &prd.Name); err != nil {
			return areas, err
		}
		areas = append(areas, prd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return areas, nil
}

// Returns the area and feature id for the given area and feature name and the product id.
func (r Repository) GetAreaAndFeatureId(area string, feature string, productId string) (int64, int64, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT a.id, f.id FROM areas a JOIN features f ON a.id = f.area_id JOIN products p ON p.id = a.product_id WHERE a.name = ? and f.name = ? and p.id=?;")
	var aid int64
	var fid int64
	if err == nil {
		defer stmt.Close()
		err = stmt.QueryRowContext(ctx, area, feature, productId).Scan(&aid, &fid)
	}
	return aid, fid, err
}
