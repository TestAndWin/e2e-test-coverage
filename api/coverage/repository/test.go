/*
Copyright (c) 2022-2026, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	component VARCHAR(255),
	url VARCHAR(500), 
	total int,  
	passes int, 
	pending int, 
	failures int, 
	skipped int, 
	testrun datetime,
	uuid VARCHAR(255),
	is_first BOOLEAN,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       FOREIGN KEY (feature_id) REFERENCES features(id),
       FOREIGN KEY (area_id) REFERENCES areas(id)
       )`

const insertTestStmt = "INSERT INTO tests (product_id, area_id, feature_id, suite, file, component, url, total, passes, pending, failures, skipped, uuid, is_first, testrun) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

const insertTestNoAreaFeatureStmt = "INSERT INTO tests (product_id, suite, file, component, url, total, passes, pending, failures, skipped, uuid, is_first, testrun) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"

const deleteTestStmt = "DELETE FROM tests WHERE component = ? AND suite = ? AND file = ?"

const testQueryPeriodDays = 28

func (cs CoverageStore) CreateTestsTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cs.db.ExecContext(ctx, CREATE_TEST)
	if err != nil {
		log.Printf("Error %s when creating Tests DB table\n", err)
		return err
	}
	return nil
}

func (cs CoverageStore) InsertTestResult(productId string, areaId int64, featureId int64, component string, url string, isFirst bool, tr reporter.TestResult) (int64, error) {
	return cs.executeSql(insertTestStmt, productId, areaId, featureId, tr.Suite, tr.File, component, url, tr.Total, tr.Passes, tr.Pending, tr.Failures, tr.Skipped, tr.Uuid, isFirst, tr.TestRun)
}

func (cs CoverageStore) InsertTestResultWithoutAreaFeature(productId string, component string, url string, isFirst bool, tr reporter.TestResult) (int64, error) {
	return cs.executeSql(insertTestNoAreaFeatureStmt, productId, tr.Suite, tr.File, component, url, tr.Total, tr.Passes, tr.Pending, tr.Failures, tr.Skipped, tr.Uuid, isFirst, tr.TestRun)
}

func (cs CoverageStore) DeleteTest(component string, suite string, file string) (int64, error) {
	return cs.executeSql(deleteTestStmt, component, suite, file)
}

// HasTestBeenUploaded checks if a test with the given UUID has already been uploaded.
// It returns true if the test exists, false otherwise, and an error if the database operation fails.
func (cs CoverageStore) HasTestBeenUploaded(uuid string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT 1 FROM tests WHERE uuid = ? LIMIT 1"

	var exists bool
	err := cs.db.QueryRowContext(ctx, query, uuid).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		log.Printf("Error %s when query context", err)
		return false, fmt.Errorf("failed to check test existence: %w", err)
	}

	return exists, nil
}

// IsThisTheFirstUpload checks if this is the first upload for the given parameters.
// It returns true if it's the first upload, false otherwise, and an error if the database operation fails.
func (cs CoverageStore) IsThisTheFirstUpload(pid string, aid int64, fid int64, suite string, file string, component string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
			SELECT 1 FROM tests 
			WHERE product_id = ? AND area_id = ? AND feature_id = ? AND suite = ? AND file = ? AND component = ? 
			LIMIT 1
	`

	var exists bool
	err := cs.db.QueryRowContext(ctx, query, pid, aid, fid, suite, file, component).Scan(&exists)

	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if this is the first upload: %w", err)
	}

	return !exists, nil
}

// Get all tests for the specified feature id
func (cs CoverageStore) GetAllFeatureTests(fid string) ([]model.Test, error) {
	return cs.GetTests(fid, "SELECT id, product_id, area_id, feature_id, suite, file, component, url, total, passes, pending, failures, skipped, uuid, is_first, testrun FROM tests WHERE feature_id = ? AND testrun > ? ORDER BY component, suite, file, testrun DESC;")
}

// Get all tests for the specified product id
func (cs CoverageStore) GetAllProductTests(pid string) ([]model.Test, error) {
	return cs.GetTests(pid, "SELECT id, product_id, COALESCE(area_id,0) as area_id, COALESCE(feature_id,0) as feature_id, suite, file, component, url, total, passes, pending, failures, skipped, uuid, is_first, testrun FROM tests WHERE product_id = ? AND testrun > ? ORDER BY component, suite, file, testrun DESC;")
}

// GetTests retrieves tests for a given ID within the last 28 days.
func (cs CoverageStore) GetTests(id string, query string) ([]model.Test, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := cs.db.QueryContext(ctx, query, id, time.Now().AddDate(0, 0, -testQueryPeriodDays))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tests []model.Test
	var prevTest *model.Test

	for rows.Next() {
		currentTest := model.Test{}
		if err := scanTest(rows, &currentTest); err != nil {
			log.Printf("Error %s when query context", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if shouldAddNewTest(prevTest, currentTest) {
			currentTest = initializeNewTest(currentTest)
			tests = append(tests, currentTest)
		} else {
			updateExistingTest(&tests[len(tests)-1], prevTest, currentTest)
		}

		prevTest = &currentTest
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tests, nil
}

func scanTest(rows *sql.Rows, t *model.Test) error {
	return rows.Scan(&t.Id, &t.ProductId, &t.AreaId, &t.FeatureId, &t.Suite, &t.FileName, &t.Component,
		&t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.IsFirst, &t.TestRun)
}

func shouldAddNewTest(prev *model.Test, current model.Test) bool {
	return prev == nil || prev.Suite != current.Suite || prev.Component != current.Component ||
		(prev.Component == current.Component && prev.FileName != current.FileName)
}

func initializeNewTest(t model.Test) model.Test {
	t.TotalTestRuns = 1
	t.FailedTestRuns = 0
	if t.Failures > 0 {
		t.FailedTestRuns = 1
	}
	if !t.IsFirst {
		t.FirstTotal = t.Total
	}
	return t
}

func updateExistingTest(existing *model.Test, prev *model.Test, current model.Test) {
	existing.FirstTotal = existing.FirstTotal - prev.Total + current.Total
	if current.Failures > 0 {
		existing.FailedTestRuns++
	}
	existing.TotalTestRuns++
}

// Get the test coverage information for all areas of the specified procduct
func (cs CoverageStore) GetAreaCoverageForProduct(productId string) (map[int64]model.Test, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	builder := sq.Select("t.id", "t.product_id", "t.area_id", "t.feature_id", "t.suite", "t.file", "t.component", "t.url",
		"t.total", "t.passes", "t.pending", "t.failures", "t.skipped", "t.uuid", "t.is_first", "t.testrun").
		From("tests t").
		Join("areas a ON a.id = t.area_id").
		Where("a.product_id = ?", productId).
		Where("t.testrun > ?", time.Now().AddDate(0, 0, -testQueryPeriodDays)).
		OrderBy("t.area_id", "t.feature_id", "t.component", "t.suite", "t.file", "t.testrun DESC")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cs.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Map where the key is the area id
	coverage := make(map[int64]model.Test)
	var prevRow *model.Test

	for rows.Next() {
		t := model.Test{}
		if err := rows.Scan(&t.Id, &t.ProductId, &t.AreaId, &t.FeatureId, &t.Suite, &t.FileName, &t.Component, &t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.IsFirst, &t.TestRun); err != nil {
			log.Println(err)
			return nil, err
		}
		if r, ok := coverage[t.AreaId]; ok {
			// We already know the area
			if prevRow == nil || prevRow.FeatureId != t.FeatureId || prevRow.Suite != t.Suite || prevRow.Component != t.Component || (prevRow.Component == t.Component && prevRow.FileName != t.FileName) {
				// Check if it is a different feature
				// And only the latest test result is taken into account, if there is more than one for the suite and file
				r.Total += t.Total
				r.Passes += t.Passes
				r.Pending += t.Pending
				r.Failures += t.Failures
				r.Skipped += t.Skipped
				// When this test has not been executed in the last 28d for the first time, we use the total value
				if !t.IsFirst {
					r.FirstTotal += t.Total
				}
			} else {
				// A bit tricky: FirstTotal should store the number of tests starting the period.
				// But we sorted the SQL result desc, so we have to get the total value from the first test
				// And not to run another SQL, we do it this way
				// With this we first remove the number of previous test (same suite, file, ...) and then adding the current one.
				r.FirstTotal = r.FirstTotal - prevRow.Total + t.Total
			}
			coverage[t.AreaId] = r
		} else {
			// new area
			// When this test has not been executed in the last 28d for the first time, we start with the total number of tests
			if !t.IsFirst {
				t.FirstTotal = t.Total
			}
			coverage[t.AreaId] = t
		}
		prevRow = &t
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return coverage, nil
}

// Get the test coverage information for all features of the specified area
func (cs CoverageStore) GetFeatureCoverageForArea(areaId string) (map[int64]model.Test, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	builder := sq.Select("t.id", "t.product_id", "t.area_id", "t.feature_id", "t.suite", "t.file", "t.component", "t.url",
		"t.total", "t.passes", "t.pending", "t.failures", "t.skipped", "t.uuid", "t.is_first", "t.testrun").
		From("tests t").
		Where("t.area_id = ?", areaId).
		Where("t.testrun > ?", time.Now().AddDate(0, 0, -testQueryPeriodDays)).
		OrderBy("t.feature_id", "t.component", "t.suite", "t.file", "t.testrun DESC")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cs.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	coverage := make(map[int64]model.Test)
	var prevRow *model.Test

	for rows.Next() {
		t := model.Test{}
		if err := rows.Scan(&t.Id, &t.ProductId, &t.AreaId, &t.FeatureId, &t.Suite, &t.FileName, &t.Component, &t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.IsFirst, &t.TestRun); err != nil {
			log.Println(err)
			return nil, err
		}

		if r, ok := coverage[t.FeatureId]; ok {
			if prevRow == nil || prevRow.Suite != t.Suite || prevRow.Component != t.Component || (prevRow.Component == t.Component && prevRow.FileName != t.FileName) {
				// Only the latest test result is taken into account, if there more than one for suite and file
				r.Total += t.Total
				r.Passes += t.Passes
				r.Pending += t.Pending
				r.Failures += t.Failures
				r.Skipped += t.Skipped
				// When this test has not been executed in the last 28d for the first time, we use the total value
				if !t.IsFirst {
					r.FirstTotal += t.Total
				}
			} else {
				// see GetAreaCoverageForProduct for explanation
				r.FirstTotal = r.FirstTotal - prevRow.Total + t.Total
			}
			coverage[t.FeatureId] = r
		} else {
			// When this test has not been executed in the last 28d for the first time, we start with the total number of tests
			if !t.IsFirst {
				t.FirstTotal = t.Total
			}
			coverage[t.FeatureId] = t
		}
		prevRow = &t
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return coverage, nil
}

// Get all tests for the specified suite and file
func (cs CoverageStore) GetAllTestForSuiteFile(component string, suite string, file string) ([]model.Test, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	builder := sq.Select("id", "product_id", "suite", "file", "component", "url", "total", "passes", "pending",
		"failures", "skipped", "uuid", "is_first", "testrun").
		From("tests").
		Where("component = ?", component).
		Where("suite = ?", suite).
		Where("file = ?", file).
		Where("testrun > ?", time.Now().AddDate(0, 0, -testQueryPeriodDays)).
		OrderBy("testrun DESC")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := cs.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error %s when querying context", err)
		return nil, err
	}

	defer rows.Close()
	var tests = []model.Test{}
	for rows.Next() {
		t := model.Test{}
		if err := rows.Scan(&t.Id, &t.ProductId, &t.Suite, &t.FileName, &t.Component, &t.Url, &t.Total, &t.Passes, &t.Pending, &t.Failures, &t.Skipped, &t.Uuid, &t.IsFirst, &t.TestRun); err != nil {
			log.Printf("Error %s when query context", err)
			return tests, err
		}
		tests = append(tests, t)
	}
	if err := rows.Err(); err != nil {
		return tests, err
	}
	return tests, nil
}

// GetComponents retrieves all components with their latest test run statistics.
func (cs CoverageStore) GetComponents() ([]model.Component, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	subquery := sq.Select("component", "MAX(testrun) AS testrun").
		From("tests").
		GroupBy("component")
	builder := sq.Select("c.component", "c.testrun",
		"SUM(t.total) as total", "SUM(t.passes) as passes",
		"SUM(t.pending) as pending", "SUM(t.failures) as failures",
		"SUM(t.skipped) as skipped").
		FromSelect(subquery, "c").
		Join("tests t ON c.component = t.component AND c.testrun = t.testrun").
		GroupBy("c.component", "c.testrun").
		OrderBy("c.component")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cs.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var components []model.Component
	for rows.Next() {
		var c model.Component
		err := rows.Scan(&c.Name, &c.TestRun, &c.Total, &c.Passes, &c.Pending, &c.Failures, &c.Skipped)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		components = append(components, c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error %s when query context", err)
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return components, nil
}

// DeleteTestsByFeatureId removes all tests associated with a specific feature
const deleteTestsByFeatureIdStmt = "DELETE FROM tests WHERE feature_id = ?"

func (cs CoverageStore) DeleteTestsByFeatureId(featureId string) (int64, error) {
	log.Printf("Deleting all tests for feature ID: %s", featureId)
	return cs.executeSql(deleteTestsByFeatureIdStmt, featureId)
}
