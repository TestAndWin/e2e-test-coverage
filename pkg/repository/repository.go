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

// Inserts/Deletes a row using the specified statement and params
func (r Repository) executeSql(statement string, params ...any) (int64, error) {
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
