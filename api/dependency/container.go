/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dependency

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/TestAndWin/e2e-coverage/auth"
	"github.com/TestAndWin/e2e-coverage/config"
	"github.com/TestAndWin/e2e-coverage/coverage/repository"
	"github.com/TestAndWin/e2e-coverage/db"
	"github.com/TestAndWin/e2e-coverage/errors"
	userRepo "github.com/TestAndWin/e2e-coverage/user/repository"
)

// Container manages all application dependencies
type Container struct {
	// Database connections
	coverageDB *sql.DB
	userDB     *sql.DB

	// Repositories
	coverageStore  *repository.CoverageStore
	userStore      *userRepo.UserStore
	
	// Auth
	tokenManager   *auth.TokenManager
	
	// Config
	appConfig     *config.Config
	
	// Mutex for thread safety
	mu sync.Mutex
}

// Create singleton instance
var (
	instance *Container
	once     sync.Once
)

// GetContainer returns the singleton container instance
func GetContainer() *Container {
	once.Do(func() {
		instance = &Container{}
		err := instance.initialize()
		if err != nil {
			log.Fatalf("Failed to initialize dependency container: %v", err)
		}
	})
	return instance
}

// initialize sets up all dependencies
func (c *Container) initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	c.appConfig = &cfg
	
	return nil
}

// GetCoverageDB returns the coverage database connection
func (c *Container) GetCoverageDB() (*sql.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.coverageDB == nil {
		var err error
		c.coverageDB, err = db.OpenDbConnection("e2ecoverage")
		if err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("failed to open coverage DB connection: %w", err))
		}
	}
	
	return c.coverageDB, nil
}

// GetUserDB returns the user database connection
func (c *Container) GetUserDB() (*sql.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.userDB == nil {
		var err error
		c.userDB, err = db.OpenDbConnection("user")
		if err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("failed to open user DB connection: %w", err))
		}
	}
	
	return c.userDB, nil
}

// GetCoverageStore returns the coverage repository
func (c *Container) GetCoverageStore() (*repository.CoverageStore, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.coverageStore == nil {
		var err error
		c.coverageStore, err = repository.NewCoverageStore()
		if err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("failed to create coverage store: %w", err))
		}
		
		// Initialize tables
		if err := c.coverageStore.CreateAllTables(); err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("failed to create tables: %w", err))
		}
	}
	
	return c.coverageStore, nil
}

// GetUserStore returns the user repository
func (c *Container) GetUserStore() (*userRepo.UserStore, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.userStore == nil {
		var err error
		c.userStore, err = userRepo.NewUserStore()
		if err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("failed to create user store: %w", err))
		}
		
		// Initialize tables if needed
		if err := c.userStore.CreateUserTable(); err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("failed to create user table: %w", err))
		}
	}
	
	return c.userStore, nil
}

// GetTokenManager returns the authentication token manager
func (c *Container) GetTokenManager() (*auth.TokenManager, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.tokenManager == nil {
		var err error
		c.tokenManager, err = auth.NewTokenManager()
		if err != nil {
			return nil, errors.NewInternalError(fmt.Errorf("failed to create token manager: %w", err))
		}
	}
	
	return c.tokenManager, nil
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() *config.Config {
	return c.appConfig
}

// CloseConnections cleanly closes all database connections
func (c *Container) CloseConnections() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.coverageDB != nil {
		c.coverageDB.Close()
		c.coverageDB = nil
	}
	
	if c.userDB != nil {
		c.userDB.Close()
		c.userDB = nil
	}
}