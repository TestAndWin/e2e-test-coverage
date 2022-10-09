package controller

import (
	"log"
	"os"

	"github.com/TestAndWin/e2e-coverage/pkg/repository"
)

var repo = initRepository()

// Set-up the db connection and create the db tables if needed
func initRepository() repository.Repository {
	repository, err := repository.OpenDbConnection()
	if err != nil {
		log.Fatal("Error connecting to DB")
		os.Exit(1)
	}
	err = repository.CreateTables()
	if err != nil {
		log.Fatal("Error connecting to DB")
		os.Exit(1)
	}
	return repository
}
