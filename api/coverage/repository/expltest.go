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

const createExplTestStmt = `CREATE TABLE IF NOT EXISTS expl_tests (
	id INT AUTO_INCREMENT PRIMARY KEY,
	area_id INT, 
	summary TEXT,
	rating INT,
	testrun datetime,
	tester INT,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (area_id) REFERENCES areas(id)
	)`

const insertExplTestStmt = "INSERT INTO expl_tests (area_id, summary, rating, testrun, tester) VALUES (?,?,?,?,?);"

const deleteExplTestStmt = "DELETE FROM expl_tests WHERE id = ?"

func (cs CoverageStore) CreateExplTestsTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cs.db.ExecContext(ctx, createExplTestStmt)
	if err != nil {
		log.Printf("Error %s when creating Expl Tests DB table\n", err)
		return err
	}
	return nil
}

func (cs CoverageStore) InsertExplTest(et model.ExplTest) (int64, error) {
	return cs.executeSql(insertExplTestStmt, et.AreaId, et.Summary, et.Rating, et.TestRun, et.Tester)
}

func (cs CoverageStore) DeleteExplTest(id string) (int64, error) {
	return cs.executeSql(deleteExplTestStmt, id)
}

// Get all exploratory tests for the specified area
func (cs CoverageStore) GetExplTests(aid string) ([]model.ExplTest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := cs.db.QueryContext(ctx, "SELECT id, area_id, summary, rating, testrun, tester FROM expl_tests WHERE area_id = ? AND testrun > ?;", aid, time.Now().AddDate(0, 0, -28))
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var et = []model.ExplTest{}
	for rows.Next() {
		e := model.ExplTest{}
		if err := rows.Scan(&e.Id, &e.AreaId, &e.Summary, &e.Rating, &e.TestRun, &e.Tester); err != nil {
			log.Println(err)
			return et, err
		}
		et = append(et, e)
	}
	if err := rows.Err(); err != nil {
		return et, err
	}
	return et, nil
}

// Returns the number of expl. tests and the average rating for the specified area id
func (cs CoverageStore) GetExplTestOverviewForArea(aid int64) (model.Area, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var a = model.Area{}
	err := cs.db.QueryRowContext(ctx, "SELECT count(*), COALESCE(AVG(rating),0) FROM expl_tests WHERE area_id = ? AND testrun > ?;", aid, time.Now().AddDate(0, 0, -28)).Scan(&a.ExplTests, &a.ExplRating)
	return a, err
}
