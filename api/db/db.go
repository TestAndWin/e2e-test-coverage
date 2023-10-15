/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
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
	_ "github.com/lib/pq"
)

func dsnMySQL(config config.Config, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", config.DBUser, config.DBPassword, config.DBHost, dbName)
}

func dsnPostgreSQL(config config.Config, dbName string) string {
	// "postgres://%s:%s@%s/%s?sslmode=verify-full"
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, dbName)
}

// Connects to the database
func OpenDbConnection(dbName string) (*sql.DB, error) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Print("Cannot load config ", err)
		os.Exit(0)
	}

	if config.DBEngine == "mysql" {
		return getMySQLConnection(config, dbName)
	} else if config.DBEngine == "postgresql" {
		return getPostgreSQLConnection(config, dbName)
	} else {
		log.Print("Unknown DB Engine ", config.DBEngine)
		os.Exit(0)
	}

	return nil, nil
}

func getMySQLConnection(config config.Config, dbName string) (*sql.DB, error) {
	log.Println("Try to connect to MySQL database")
	db, err := sql.Open("mysql", dsnMySQL(config, ""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return db, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
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

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return db, err
	}
	log.Printf("Connected to DB %s successfully\n", dbName)
	return db, err
}

func getPostgreSQLConnection(config config.Config, dbName string) (*sql.DB, error) {
	log.Println("Try to connect to PostgreSQL database")
	db, err := sql.Open("postgres", dsnPostgreSQL(config, ""))

	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return db, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	rows, err := db.Query(fmt.Sprintf("SELECT FROM pg_database WHERE datname = '%s';", dbName))
	if err != nil {
		log.Printf("Error %s when running sql statement\n", err)
		return db, err
	}
	defer rows.Close()
	if !rows.Next() {
		_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s;", dbName))
		if err != nil {
			log.Printf("Error %s when creating DB\n", err)
			return db, err
		}

		db.Close()
		db, err = sql.Open("postgres", dsnPostgreSQL(config, dbName))
		if err != nil {
			log.Printf("Error %s when opening DB", err)
			return db, err
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(time.Minute * 5)

		ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		err = db.PingContext(ctx)
		if err != nil {
			log.Printf("Errors %s pinging DB", err)
			return db, err
		}
		log.Printf("Connected to DB %s successfully\n", dbName)
	}

	return db, err
}
