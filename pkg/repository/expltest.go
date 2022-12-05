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

	"github.com/TestAndWin/e2e-coverage/pkg/model"
)

const createExplTestStmt = `CREATE TABLE IF NOT EXISTS expl_tests (
	id INT AUTO_INCREMENT PRIMARY KEY,
	area_id INT, 
	summary TEXT,
	rating INT,
	testrun datetime,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (area_id) REFERENCES areas(id)
	)`

const insertExplTestStmt = "INSERT INTO expl_tests (area_id, summary, rating, testrun) VALUES (?,?,?,?);"

const deleteExplTestStmt = "DELETE FROM expl_tests WHERE id = ?"

func (r Repository) CreateExplTestsTable() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := r.db.ExecContext(ctx, createExplTestStmt)
	if err != nil {
		log.Printf("Error %s when creating Expl Tests DB table\n", err)
		return err
	}
	return nil
}

func (r Repository) InsertExplTest(et model.ExplTest) (int64, error) {
	return r.executeSql(insertExplTestStmt, et.AreaId, et.Summary, et.Rating, et.TestRun)
}

func (r Repository) DeleteExplTest(id string) (int64, error) {
	return r.executeSql(deleteExplTestStmt, id)
}

// Get all exploratory tests for the specified area
func (r Repository) GetExplTests(aid string) ([]model.ExplTest, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, area_id, summary, rating, testrun FROM expl_tests WHERE area_id = ? AND testrun > ?;")
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, aid, time.Now().AddDate(0, 0, -28))
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var et = []model.ExplTest{}
	for rows.Next() {
		var e model.ExplTest
		if err := rows.Scan(&e.Id, &e.AreaId, &e.Summary, &e.Rating, &e.TestRun); err != nil {
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
func (r Repository) GetExplTestOverviewForArea(aid int64) (model.Area, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT count(*), COALESCE(AVG(rating),0) FROM expl_tests WHERE area_id = ? AND testrun > ?;")
	var a = model.Area{}
	if err == nil {
		defer stmt.Close()
		err = stmt.QueryRowContext(ctx, aid, time.Now().AddDate(0, 0, -28)).Scan(&a.ExplTests, &a.ExplRating)
	}
	return a, err
}
