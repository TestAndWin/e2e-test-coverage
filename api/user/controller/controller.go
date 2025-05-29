package controller

import (
	"github.com/TestAndWin/e2e-coverage/dependency"
	"github.com/TestAndWin/e2e-coverage/user/repository"
)

// getUserRepository returns the user repository from the dependency container
func getUserRepository() (*repository.UserStore, error) {
	container := dependency.GetContainer()
	store, err := container.GetUserStore()
	if err != nil {
		return nil, err
	}
	return store, nil
}
