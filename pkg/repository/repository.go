/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
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
	"os"
	"time"

	"github.com/TestAndWin/e2e-coverage/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
}

func loadConfig() (config Config, err error) {
	viper.SetConfigFile("db.env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
		return
	}
	err = viper.Unmarshal(&config)
	return
}

type Repository struct {
	db *sql.DB
}

func dsn(config Config, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", config.DBUser, config.DBPassword, config.DBHost, dbName)
}

// Connects to the database
func OpenDbConnection() (Repository, error) {
	log.Println("Try to connect to database")
	config, err := loadConfig()
	if err != nil {
		log.Print("Cannot load config ", err)
		os.Exit(0)
	}

	db, err := sql.Open("mysql", dsn(config, ""))
	r := Repository{db: db}
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return r, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err = r.db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+config.DBName)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return r, err
	}

	r.db.Close()
	r.db, err = sql.Open("mysql", dsn(config, config.DBName))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return r, err
	}

	r.db.SetMaxOpenConns(20)
	r.db.SetMaxIdleConns(20)
	r.db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = r.db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return r, err
	}
	log.Printf("Connected to DB %s successfully\n", config.DBName)
	return r, err
}

// Creates the db tables if not exiting
func (r Repository) CreateTables() error {

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := r.db.ExecContext(ctx, model.CREATE_PRODUCT)
	if err != nil {
		log.Printf("Error %s when creating Products DB table\n", err)
		return err
	}
	_, err = r.db.ExecContext(ctx, model.CREATE_AREA)
	if err != nil {
		log.Printf("Error %s when creating Areas DB table\n", err)
		return err
	}
	_, err = r.db.ExecContext(ctx, model.CREATE_FEATURE)
	if err != nil {
		log.Printf("Error %s when creating Features DB table\n", err)
		return err
	}
	_, err = r.db.ExecContext(ctx, model.CREATE_TEST)
	if err != nil {
		log.Printf("Error %s when creating Tests DB table\n", err)
		return err
	}

	_, err = r.db.ExecContext(ctx, model.CREATE_EXPL_TEST)
	if err != nil {
		log.Printf("Error %s when creating Tests DB table\n", err)
		return err
	}
	return nil
}

// Inserts/Deletes a row using the specified statement and params
func (r Repository) ExecuteSql(statement string, params ...any) (int64, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, statement)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return 0, err
	}
	log.Printf("%d row inserted/deleted/updated ", rows)
	return res.LastInsertId()
}

// Checks if the test had already been uploaded.
func (r Repository) HasTestBeenUploaded(uuid string) (bool, error) {
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

// Returns all products
func (r Repository) GetAllProducts() ([]model.Product, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, model.SELECT_PRODUCTS)
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

// Get all areas for the specified product id
func (r Repository) GetAllProductAreas(pid string) ([]model.Area, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, model.SELECT_AREAS)
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

// Get all features for the specified area id
func (r Repository) GetAllAreaFeatures(aid string) ([]model.Feature, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, model.SELECT_FEATURES)
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

// Get all tests for the specified feature id
func (r Repository) GetAllFeatureTests(fid string) ([]model.Test, error) {
	return getTests(r, fid, "SELECT id, product_id, area_id, feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun FROM tests WHERE feature_id = ? AND testrun > ? ORDER BY suite, file, testrun DESC;")
}

// Get all tests for the specified product id
func (r Repository) GetAllProductTests(pid string) ([]model.Test, error) {
	return getTests(r, pid, "SELECT id, product_id, COALESCE(area_id,0) as area_id, COALESCE(feature_id,0) as feature_id, suite, file, url, total, passes, pending, failures, skipped, uuid, testrun FROM tests WHERE product_id = ? AND testrun > ? ORDER BY suite, file, testrun DESC;")
}

func getTests(r Repository, id string, query string) ([]model.Test, error) {
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
func (r Repository) GetAreaCoverageForProduct(productId string) (map[int64]model.Test, error) {
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
func (r Repository) GetFeatureCoverageForArea(areaId string) (map[int64]model.Test, error) {
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

// Get all exploratory tests for the specified area
func (r Repository) GetExplTests(aid string) ([]model.ExplTest, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, model.SELECT_EXPL_TESTS_28D)
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

// Get all tests for the specified suite and file
func (r Repository) GetAllTestForSuiteFile(suite string, file string) ([]model.Test, error) {
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
