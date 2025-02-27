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
	"fmt"
	"log"
	"time"

	"github.com/TestAndWin/e2e-coverage/db"
	_ "github.com/go-sql-driver/mysql"
)

// CoverageStore handles all database operations for coverage data
type CoverageStore struct {
	db *sql.DB
}

// Interface for CoverageStore to enable mocking in tests
type CoverageRepository interface {
	// Tables
	CreateProductsTable() error
	CreateAreasTable() error
	CreateExplTestsTable() error
	CreateFeaturesTable() error
	CreateTestsTable() error
	CreateAllTables() error
}

// NewCoverageStore creates a new CoverageStore instance
func NewCoverageStore() (*CoverageStore, error) {
	database, err := db.OpenDbConnection("e2ecoverage")
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	log.Println("Connected successfully to coverage database")
	return &CoverageStore{
		db: database,
	}, nil
}

// WithDB creates a new CoverageStore with a provided DB connection
func WithDB(database *sql.DB) *CoverageStore {
	return &CoverageStore{
		db: database,
	}
}

// CreateAllTables creates all required tables for the application
func (store *CoverageStore) CreateAllTables() error {
	tables := []struct {
		name string
		fn   func() error
	}{
		{"Products", store.CreateProductsTable},
		{"Areas", store.CreateAreasTable},
		{"ExplTests", store.CreateExplTestsTable},
		{"Features", store.CreateFeaturesTable},
		{"Tests", store.CreateTestsTable},
	}

	for _, table := range tables {
		if err := table.fn(); err != nil {
			return fmt.Errorf("failed to create %s table: %w", table.name, err)
		}
	}
	
	return nil
}

// Inserts/Deletes a row using the specified statement and params
func (cs CoverageStore) executeSql(statement string, params ...any) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	stmt, err := cs.db.PrepareContext(ctx, statement)
	if err != nil {
		return 0, fmt.Errorf("error preparing SQL statement: %w", err)
	}
	defer stmt.Close()
	
	res, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		return 0, fmt.Errorf("error executing statement: %w", err)
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error finding rows affected: %w", err)
	}
	
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}
	
	log.Printf("%d row(s) affected, last ID: %d", rows, lastID)
	return lastID, nil
}
