package controller

import (
	"github.com/TestAndWin/e2e-coverage/dependency"
	"github.com/TestAndWin/e2e-coverage/user/repository"
)

// getUserRepository returns the user repository from the dependency container
func getUserRepository() *repository.UserStore {
	container := dependency.GetContainer()
	store, err := container.GetUserStore()
	if err != nil {
		panic(err) // This should never happen as the container handles initialization errors
	}
	return store
}
