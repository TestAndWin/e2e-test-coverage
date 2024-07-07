/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TestAndWin/e2e-coverage/config"
	_ "github.com/go-sql-driver/mysql"
)

func dsnMySQL(config config.Config, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", config.DBUser, config.DBPassword, config.DBHost, dbName)
}

// Connects to the database
func OpenDbConnection(dbName string) (*sql.DB, error) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Print("Cannot load config ", err)
		os.Exit(0)
	}

	log.Println("Try to connect to MySQL database")
	db, err := sql.Open("mysql", dsnMySQL(config, ""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return db, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbName)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return db, err
	}

	db.Close()
	db, err = sql.Open("mysql", dsnMySQL(config, dbName))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return db, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return db, err
	}
	log.Printf("Connected to DB %s successfully\n", dbName)
	return db, err
}
