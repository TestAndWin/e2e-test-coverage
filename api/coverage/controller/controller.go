/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"log"
	"os"

	"github.com/TestAndWin/e2e-coverage/coverage/repository"
	"github.com/gin-gonic/gin"
)

var repo = initRepository()

// Set-up the db connection and create the db tables if needed
func initRepository() *repository.CoverageStore {
	repository, err := repository.OpenDbConnection()
	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err)
		os.Exit(1)
	}

	err = repository.CreateProductsTable()
	if err != nil {
		log.Fatalf("Error creating DB tables: %s", err)
		os.Exit(1)
	}
	err = repository.CreateAreasTable()
	if err != nil {
		log.Fatalf("Error creating DB tables: %s", err)
		os.Exit(1)
	}
	err = repository.CreateExplTestsTable()
	if err != nil {
		log.Fatalf("Error creating DB tables: %s", err)
		os.Exit(1)
	}
	err = repository.CreateFeaturesTable()
	if err != nil {
		log.Fatalf("Error creating DB tables: %s", err)
		os.Exit(1)
	}
	err = repository.CreateTestsTable()
	if err != nil {
		log.Fatalf("Error creating DB tables: %s", err)
		os.Exit(1)
	}

	return repository
}

func handleError(c *gin.Context, err error, message string, status int) {
	log.Printf("%s: %v", message, err)
	c.JSON(status, gin.H{
		"error":   message,
		"details": err.Error(),
		"status":  status,
	})
}
