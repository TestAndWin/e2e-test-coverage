/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
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

	"github.com/TestAndWin/e2e-coverage/db"
	_ "github.com/go-sql-driver/mysql"
)

type CoverageStore struct {
	db *sql.DB
}

// Creates a new CoverageStore connecting to a MySQL database
func OpenDbConnection() (*CoverageStore, error) {
	db, err := db.OpenDbConnection("e2ecoverage")
	if err != nil {
		return nil, err
	}
	return &CoverageStore{db: db}, nil
}

// Inserts/Deletes a row using the specified statement and params
func (r CoverageStore) executeSql(statement string, params ...any) (int64, error) {
	//log.Printf("statement %s, params: %s", statement, params)
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
		log.Printf("Error %s when executing statement", err)
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
