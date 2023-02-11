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
	"github.com/TestAndWin/e2e-coverage/coverage/reporter"
)

const CREATE_TEST = `CREATE TABLE IF NOT EXISTS tests (
	id INT AUTO_INCREMENT PRIMARY KEY,
	product_id int,
	area_id int,
	feature_id int,
	suite VARCHAR(255),  
	file VARCHAR(255),
	url VARCHAR(500), 
	total int,  
	passes int, 
	pending int, 
	failures int, 
	skipped int, 
	testrun datetime,
	uuid VARCHAR(255),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (feature_id) REFERENCES features(id),
	FOREIGN KEY (feature_id) REFERENCES features(id)
	)`

const insertTestStmt = "INSERT INTO tests (product_id, area_id, feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"

const insertTestNoAreaFeatureStmt = "INSERT INTO tests (product_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun) VALUES (?,?,?,?,?,?,?,?,?,?,?)"

const deleteTestStmt = "DELETE FROM tests WHERE id = ?"

func (r CoverageStore) CreateTestsTable() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := r.db.ExecContext(ctx, CREATE_TEST)
	if err != nil {
		log.Printf("Error %s when creating Tests DB table\n", err)
		return err
	}
	return nil
}

func (r CoverageStore) InsertTest(t model.Test) (int64, error) {
	return r.executeSql(insertTestStmt, t.ProductId, t.AreaId, t.FeatureId, t.Suite, t.FileName, t.Url, t.Total, t.Passes, t.Pending, t.Failures, t.Skipped, t.Uuid, t.TestRun)
}

func (r CoverageStore) InsertTestResult(productId string, areadId int64, featureId int64, url string, tr reporter.TestResult) (int64, error) {
	return r.executeSql(insertTestStmt, productId, areadId, featureId, tr.Suite, tr.File, url, tr.Total, tr.Passes, tr.Pending, tr.Failures, tr.Skipped, tr.Uuid, tr.TestRun)
}

func (r CoverageStore) InsertTestResultWithoutAreaFeature(productId string, url string, tr reporter.TestResult) (int64, error) {
	return r.executeSql(insertTestNoAreaFeatureStmt, productId, tr.Suite, tr.File, url, tr.Total, tr.Passes, tr.Pending, tr.Failures, tr.Skipped, tr.Uuid, tr.TestRun)
}

func (r CoverageStore) DeleteTest(id string) (int64, error) {
	return r.executeSql(deleteTestStmt, id)
}

// Checks if the test had already been uploaded.
func (r CoverageStore) HasTestBeenUploaded(uuid string) (bool, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT count(*) FROM tests WHERE uuid = ?;")
	var count int64
	if err == nil {
		defer stmt.Close()
		err = stmt.QueryRowContext(ctx, uuid).Scan(&count)
	}
	return count > 0, err
}

// Get all tests for the specified feature id
func (r CoverageStore) GetAllFeatureTests(fid string) ([]model.Test, error) {
	return getTests(r, fid, "SELECT id, product_id, area_id, feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun FROM tests WHERE feature_id = ? AND testrun > ? ORDER BY suite, file, testrun DESC;")
}

// Get all tests for the specified product id
func (r CoverageStore) GetAllProductTests(pid string) ([]model.Test, error) {
	return getTests(r, pid, "SELECT id, product_id, COALESCE(area_id,0) as area_id, COALESCE(feature_id,0) as feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun FROM tests WHERE product_id = ? AND testrun > ? ORDER BY suite, file, testrun DESC;")
}

func getTests(r CoverageStore, id string, query string) ([]model.Test, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id, time.Now().AddDate(0, 0, -28))
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var tests = []model.Test{}
	var prevRow model.Test = model.Test{}
	for rows.Next() {
		var t model.Test
		if err := rows.Scan(&t.Id, &t.ProductId, &t.AreaId, &t.FeatureId, &t.Suite, &t.FileName, &t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.TestRun); err != nil {
			return tests, err
		}
		if prevRow.Suite != t.Suite || (prevRow.Suite == t.Suite && prevRow.FileName != t.FileName) {
			t.TotalTestRuns = 1

			t.FailedTestRuns = 0
			if t.Failures > 0 {
				t.FailedTestRuns = 1
			}

			t.FirstTotal = t.Total
			tests = append(tests, t)
		} else {
			p := tests[len(tests)-1]
			p.FirstTotal = p.FirstTotal - prevRow.Total + t.Total
			if t.Failures > 0 {
				p.FailedTestRuns += 1
			}
			p.TotalTestRuns += 1
			tests[len(tests)-1] = p
		}
		prevRow = t
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

// Get the test coverage information for all areas of the specified procduct
func (r CoverageStore) GetAreaCoverageForProduct(productId string) (map[int64]model.Test, error) {
	statement := "SELECT t.id, t.product_id, t.area_id, t.feature_id, t.suite, t.file, t.url, t.total, t.passes, t.pending, t.failures, t.skipped, t.uuid, t.testrun FROM tests t JOIN areas a ON a.id = t.area_id WHERE a.product_id = ? and t.testrun > ? ORDER BY t.area_id, t.feature_id, t.suite, t.file, t.testrun DESC;"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, statement)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()

	// Only test from the last 28 days
	rows, err := stmt.QueryContext(ctx, productId, time.Now().AddDate(0, 0, -28))
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}
	defer rows.Close()
	var coverage map[int64]model.Test = make(map[int64]model.Test)
	var prevRow model.Test = model.Test{}

	for rows.Next() {
		var t model.Test
		if err := rows.Scan(&t.Id, &t.ProductId, &t.AreaId, &t.FeatureId, &t.Suite, &t.FileName, &t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.TestRun); err != nil {
			log.Println(err)
			return nil, err
		}
		if r, ok := coverage[t.AreaId]; ok {
			// We already know the area
			if prevRow.FeatureId != t.FeatureId || prevRow.Suite != t.Suite || (prevRow.Suite == t.Suite && prevRow.FileName != t.FileName) {
				// Only the latest test result is taken into account, if there is more than one for the suite and file
				r.Total += t.Total
				r.Passes += t.Passes
				r.Pending += t.Pending
				r.Failures += t.Failures
				r.Skipped += t.Skipped
				r.FirstTotal += t.Total
			} else {
				// A bit tricky: FirstTotal stores the number of tests starting the period.
				// And not to run another SQL, we must store the value or the first test.
				// With this we first remove the number of previous test (same suite, file, ...) and then adding the current one.
				r.FirstTotal = r.FirstTotal - prevRow.Total + t.Total
			}
			coverage[t.AreaId] = r
		} else {
			// new area
			t.FirstTotal = t.Total
			coverage[t.AreaId] = t
		}
		prevRow = t
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return coverage, nil
}

// Get the test coverage information for all features of the specified area
func (r CoverageStore) GetFeatureCoverageForArea(areaId string) (map[int64]model.Test, error) {
	statement := "SELECT t.id, t.product_id, t.area_id, t.feature_id, t.suite, t.file, t.url, t.total, t.passes, t.pending, t.failures, t.skipped, t.uuid, t.testrun FROM tests t WHERE t.area_id = ? and t.testrun > ? ORDER BY t.feature_id, t.suite, t.file, t.testrun DESC;"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, statement)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()

	// Only test from the last 28 days
	rows, err := stmt.QueryContext(ctx, areaId, time.Now().AddDate(0, 0, -28))
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}
	defer rows.Close()
	var coverage map[int64]model.Test = make(map[int64]model.Test)
	var prevRow model.Test = model.Test{}

	for rows.Next() {
		var t model.Test
		if err := rows.Scan(&t.Id, &t.ProductId, &t.AreaId, &t.FeatureId, &t.Suite, &t.FileName, &t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.TestRun); err != nil {
			log.Println(err)
			return nil, err
		}

		if r, ok := coverage[t.FeatureId]; ok {
			if prevRow.Suite != t.Suite || prevRow.Suite == t.Suite && prevRow.FileName != t.FileName {
				// Only the latest test result is taken into account, if there more than one for suite and file
				r.Total += t.Total
				r.Passes += t.Passes
				r.Pending += t.Pending
				r.Failures += t.Failures
				r.Skipped += t.Skipped
				r.FirstTotal += t.Total
			} else {
				// see GetAreaCoverageForProduct for explanation
				r.FirstTotal = r.FirstTotal - prevRow.Total + t.Total
			}
			coverage[t.FeatureId] = r
		} else {
			t.FirstTotal = t.Total
			coverage[t.FeatureId] = t
		}
		prevRow = t
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return coverage, nil
}

// Get all tests for the specified suite and file
func (r CoverageStore) GetAllTestForSuiteFile(suite string, file string) ([]model.Test, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, product_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun FROM tests WHERE suite = ? AND file = ? ORDER BY testrun DESC;")
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, suite, file)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var tests = []model.Test{}
	for rows.Next() {
		var t model.Test
		if err := rows.Scan(&t.Id, &t.ProductId, &t.Suite, &t.FileName, &t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.TestRun); err != nil {
			return tests, err
		}
		tests = append(tests, t)
	}
	if err := rows.Err(); err != nil {
		return tests, err
	}
	return tests, nil
}
