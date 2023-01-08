package controller

import (
	"log"
	"os"

	"github.com/TestAndWin/e2e-coverage/user/repository"
)

var userStore = initRepository()

// Set-up the db connection and create the db tables if needed
func initRepository() repository.UserStore {
	userStore, err := repository.NewUserStore()
	if err != nil {
		log.Fatal("Error connecting to DB")
		os.Exit(1)
	}

	err = userStore.CreateUsersTable()
	if err != nil {
		log.Fatal("Error creating tables")
		os.Exit(1)
	}
	return *userStore
}
