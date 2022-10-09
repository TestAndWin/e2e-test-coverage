package main

import (
	"fmt"

	_ "github.com/TestAndWin/e2e-coverage/docs"

	"github.com/TestAndWin/e2e-coverage/pkg/router"
)

// @title e2ecoverage
// @version 1.0
// @description API for e2e-coverage

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	fmt.Println("**** e2e test coverage ****")

	router.HandleRequest()
}
