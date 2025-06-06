/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func (cs CoverageStore) CreateAreasTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cs.db.ExecContext(ctx, createAreaStmt)
	if err != nil {
		log.Printf("Error %s when creating Areas DB table\n", err)
		return err
	}
	return nil
}

func (cs CoverageStore) InsertArea(a model.Area) (int64, error) {
	return cs.executeSql(insertAreaStmt, a.ProductId, a.Name)
}

func (cs CoverageStore) UpdateArea(a model.Area) (int64, error) {
	return cs.executeSql(updateAreaStmt, a.Name, a.Id)
}

func (cs CoverageStore) DeleteArea(id int64) (int64, error) {
	return cs.executeSql(deleteAreaStmt, id)
}

// Get all areas for the specified product id
func (cs CoverageStore) GetAllProductAreas(pid string) ([]model.Area, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := cs.db.QueryContext(ctx, "SELECT id, product_id, name FROM areas WHERE product_id = ? ORDER BY name;", pid)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var areas = []model.Area{}
	for rows.Next() {
		a := model.Area{}
		if err := rows.Scan(&a.Id, &a.ProductId, &a.Name); err != nil {
			log.Println(err)
			return areas, err
		}
		areas = append(areas, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return areas, nil
}

// Returns the area and feature id for the given area and feature name and the product id.
func (cs CoverageStore) GetAreaAndFeatureId(area string, feature string, productId string) (int64, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var aid int64
	var fid int64
	builder := sq.Select("a.id", "f.id").
		From("areas a").
		Join("features f ON a.id = f.area_id").
		Join("products p ON p.id = a.product_id").
		Where("a.name = ?", area).
		Where("f.name = ?", feature).
		Where("p.id = ?", productId)
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, 0, err
	}
	err = cs.db.QueryRowContext(ctx, query, args...).Scan(&aid, &fid)
	return aid, fid, err
}

// GetAreaIdByNameAndProductId returns the area ID for an area with the given name in the specified product.
// Returns 0 and sql.ErrNoRows if no matching area is found.
func (cs CoverageStore) GetAreaIdByNameAndProductId(areaName string, productId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var areaId int64
	err := cs.db.QueryRowContext(ctx,
		"SELECT id FROM areas WHERE name = ? AND product_id = ?",
		areaName, productId).Scan(&areaId)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		log.Printf("Error getting area ID: %v", err)
		return 0, err
	}

	return areaId, nil
}
