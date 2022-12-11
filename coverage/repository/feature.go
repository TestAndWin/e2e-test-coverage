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

const createFeatureStmt = `CREATE TABLE IF NOT EXISTS features(
	id INT AUTO_INCREMENT PRIMARY KEY,
	area_id INT, 
	name VARCHAR(255),
	documentation VARCHAR(255),
	url VARCHAR(255),
	business_value VARCHAR(20),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (area_id) REFERENCES areas(id)
	)`

const insertFeatureStmt = "INSERT INTO features(area_id, name, documentation, url, business_value) VALUES (?,?,?,?,?)"

const updateFeatureStmt = "UPDATE features SET name = ?, documentation = ?, url = ?, business_value = ? WHERE id = ?"

const deleteFeatureStmt = "DELETE FROM features WHERE id = ?"

func (r Repository) CreateFeaturesTable() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := r.db.ExecContext(ctx, createFeatureStmt)
	if err != nil {
		log.Printf("Error %s when creating Features DB table\n", err)
		return err
	}
	return nil
}

func (r Repository) InsertFeature(f model.Feature) (int64, error) {
	return r.executeSql(insertFeatureStmt, f.AreaId, f.Name, f.Documentation, f.Url, f.BusinessValue)
}

func (r Repository) UpdateFeature(f model.Feature) (int64, error) {
	return r.executeSql(updateFeatureStmt, f.Name, f.Documentation, f.Url, f.BusinessValue, f.Id)
}

func (r Repository) DeleteFeature(id string) (int64, error) {
	return r.executeSql(deleteFeatureStmt, id)
}

// Get all features for the specified area id
func (r Repository) GetAllAreaFeatures(aid string) ([]model.Feature, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, area_id, name, documentation, url, business_value FROM features WHERE area_id = ?;")
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, aid)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var features = []model.Feature{}
	for rows.Next() {
		var prd model.Feature
		if err := rows.Scan(&prd.Id, &prd.AreaId, &prd.Name, &prd.Documentation, &prd.Url, &prd.BusinessValue); err != nil {
			return features, err
		}
		features = append(features, prd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return features, nil
}
