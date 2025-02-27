/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
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

func (cs CoverageStore) CreateFeaturesTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cs.db.ExecContext(ctx, createFeatureStmt)
	if err != nil {
		log.Printf("Error %s when creating Features DB table\n", err)
		return err
	}
	return nil
}

func (cs CoverageStore) InsertFeature(f model.Feature) (int64, error) {
	return cs.executeSql(insertFeatureStmt, f.AreaId, f.Name, f.Documentation, f.Url, f.BusinessValue)
}

func (cs CoverageStore) UpdateFeature(f model.Feature) (int64, error) {
	return cs.executeSql(updateFeatureStmt, f.Name, f.Documentation, f.Url, f.BusinessValue, f.Id)
}

func (cs CoverageStore) DeleteFeature(id string) (int64, error) {
	return cs.executeSql(deleteFeatureStmt, id)
}

// Get all features for the specified area id
func (cs CoverageStore) GetAllAreaFeatures(aid string) ([]model.Feature, error) {
	log.Printf("GetAllAreaFeatures: Looking for features with area_id = %s", aid)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := cs.db.QueryContext(ctx, "SELECT id, area_id, name, documentation, url, business_value FROM features WHERE area_id = ? ORDER BY name;", aid)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var features = []model.Feature{}
	rowCount := 0
	
	for rows.Next() {
		rowCount++
		f := model.Feature{}
		if err := rows.Scan(&f.Id, &f.AreaId, &f.Name, &f.Documentation, &f.Url, &f.BusinessValue); err != nil {
			log.Printf("Error scanning feature row: %v", err)
			return features, err
		}
		log.Printf("Found feature: id=%d, area_id=%d, name=%s", f.Id, f.AreaId, f.Name)
		features = append(features, f)
	}
	
	log.Printf("GetAllAreaFeatures: Found %d features for area_id %s", rowCount, aid)
	
	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v", err)
		return nil, err
	}
	
	// Additional debug for empty result
	if len(features) == 0 {
		// Check if area exists
		var areaExists bool
		err := cs.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM areas WHERE id = ?)", aid).Scan(&areaExists)
		if err != nil {
			log.Printf("Error checking if area exists: %v", err)
		} else {
			log.Printf("Area with id %s exists: %v", aid, areaExists)
		}
		
		// Count features in the database
		var totalFeatures int
		err = cs.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM features").Scan(&totalFeatures)
		if err != nil {
			log.Printf("Error counting features: %v", err)
		} else {
			log.Printf("Total features in database: %d", totalFeatures)
		}
	}
	
	return features, nil
}
