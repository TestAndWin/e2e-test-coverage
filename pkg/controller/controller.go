package controller

import (
	"log"
	"os"

	"github.com/TestAndWin/e2e-coverage/pkg/repository"
	"github.com/gin-gonic/gin"
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func OptionMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
}
