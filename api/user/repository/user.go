/*
Copyright (c) 2022-2026, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TestAndWin/e2e-coverage/db"
	"github.com/TestAndWin/e2e-coverage/user/model"
	"golang.org/x/crypto/bcrypt"
)

const createUserTable = `CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(255) UNIQUE,
	password VARCHAR(255),
	api_key VARCHAR(100),
	role VARCHAR(50),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUserTable() error
	GetUserById(id int64) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	Login(email, password string) (model.User, error)
	GetUser() ([]model.User, error)
	CreateUser(user model.User) (int64, error)
	UpdateUser(user model.User) error
	DeleteUser(id int64) error
	ChangePassword(id int64, oldPwd, newPwd string) error
	GenerateApiKey(userId int64) (string, error)
	GetUserIdForApiKey(apiKey string) (int64, error)
}

// UserStore is a store for user data that uses a MySQL database
type UserStore struct {
	db *sql.DB
}

// NewUserStore creates a new UserStore that uses a MySQL database
func NewUserStore() (*UserStore, error) {
	database, err := db.OpenDbConnection("user")
	if err != nil {
		return nil, fmt.Errorf("failed to open user DB connection: %w", err)
	}

	log.Println("Connected successfully to user database")
	return &UserStore{db: database}, nil
}

// WithDB creates a new UserStore with a provided DB connection
func WithDB(database *sql.DB) *UserStore {
	return &UserStore{
		db: database,
	}
}

// CreateUserTable creates the users table if not existing and create an admin user
func (s *UserStore) CreateUserTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create db table
	_, err := s.db.ExecContext(ctx, createUserTable)
	if err != nil {
		log.Printf("Error %s when creating Users DB table\n", err)
		return err
	}

	// Create Admin user if not created yet
	var count int64
	err = s.db.QueryRow("SELECT count(*) FROM users;").Scan(&count)
	if err != nil {
		return err
	}
	if count < 1 {
		user := model.User{Email: "admin", Password: "e2ecoverage", Roles: []string{model.ADMIN}}
		_, err = s.CreateUser(user)
	}
	return err
}

// CreateUser creates a new user. The password will be stored encrypted. It returns the id of the new user.
func (s *UserStore) CreateUser(user model.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// We keep it simple for now
	return s.executeSql("INSERT INTO users(email, password, role) VALUES (?,?,?)", user.Email, hashedPassword, strings.Join(user.Roles, ","))
}

// UpdateUser updates a user's information
func (s *UserStore) UpdateUser(user model.User) error {
	if len(user.Email) > 0 {
		_, err := s.executeSql("UPDATE users SET email = ? WHERE id = ?", user.Email, user.Id)
		if err != nil {
			return err
		}
	}

	_, err := s.executeSql("UPDATE users SET role = ? WHERE id = ?", strings.Join(user.Roles, ","), user.Id)
	if err != nil {
		return err
	}

	if len(user.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = s.executeSql("UPDATE users SET password = ? WHERE id = ?", hashedPassword, user.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

// ChangePassword changes a user's password after validating their current password
func (s *UserStore) ChangePassword(id int64, oldPassword string, newPassword string) error {
	// Get user email for validation
	var email string
	err := s.db.QueryRow("SELECT email FROM users WHERE id = ?", id).Scan(&email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Validate old password
	_, err = s.Login(email, oldPassword)
	if err != nil {
		return fmt.Errorf("invalid current password: %w", err)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password
	_, err = s.executeSql("UPDATE users SET password = ? WHERE id = ?", hashedPassword, id)
	return err
}

// Checks if the user can login and if yes, returns the id, roles, otherwise an error
func (s *UserStore) Login(email string, password string) (model.User, error) {
	var role string
	var user = model.User{}
	user.Email = email
	var hashedPassword []byte
	err := s.db.QueryRow("SELECT id, password, role FROM users WHERE email = ?", email).Scan(&user.Id, &hashedPassword, &role)
	if err != nil {
		return user, err
	}
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return user, err
	}

	user.Roles = strings.Split(role, ",")
	return user, nil
}

// Inserts/Deletes a row using the specified statement and params
func (s *UserStore) executeSql(statement string, params ...any) (int64, error) {
	log.Printf("statement: %s, params: %v", statement, params)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, statement, params...)
	if err != nil {
		return 0, fmt.Errorf("error executing statement: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error finding rows affected: %w", err)
	}

	log.Printf("%d rows affected", rows)
	return res.LastInsertId()
}

// Returns all user
func (s *UserStore) GetUser() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, "SELECT id, email, role FROM users;")
	if err != nil {
		log.Printf("Error %s when query context", err)
		return nil, err
	}

	defer rows.Close()
	var user = []model.User{}
	for rows.Next() {
		var u model.User
		var r string
		if err := rows.Scan(&u.Id, &u.Email, &r); err != nil {
			return user, err
		}
		u.Roles = strings.Split(r, ",")
		user = append(user, u)
	}
	if err := rows.Err(); err != nil {
		return []model.User{}, err
	}
	return user, nil
}

// Delete the user
func (s *UserStore) DeleteUser(id int64) error {
	_, err := s.executeSql("DELETE FROM users WHERE id = ?", id)
	return err
}

// Create a new API Key and store it in the DB.
func (s *UserStore) GenerateApiKey(userId int64) (string, error) {
	key := fmt.Sprintf("%d%s", userId, generateAPIKey())
	_, err := s.executeSql("UPDATE users SET api_key = ? WHERE id = ?", key, userId)
	return key, err
}

// Get the userId for the specified apiKey
func (s *UserStore) GetUserIdForApiKey(apiKey string) (int64, error) {
	var userId int64
	err := s.db.QueryRow("SELECT id FROM users WHERE api_key = ?", apiKey).Scan(&userId)
	return userId, err
}

func generateAPIKey() string {
	size := 32
	key := make([]byte, size)
	rand.Read(key)
	hash := sha256.Sum256(key)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// GetUserById retrieves a user by ID
func (s *UserStore) GetUserById(id int64) (model.User, error) {
	var user model.User
	var role string

	err := s.db.QueryRow("SELECT id, email, role FROM users WHERE id = ?", id).Scan(&user.Id, &user.Email, &role)
	if err != nil {
		return user, fmt.Errorf("failed to get user by id: %w", err)
	}

	user.Roles = strings.Split(role, ",")
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserStore) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	var role string

	err := s.db.QueryRow("SELECT id, email, role FROM users WHERE email = ?", email).Scan(&user.Id, &user.Email, &role)
	if err != nil {
		return user, fmt.Errorf("failed to get user by email: %w", err)
	}

	user.Roles = strings.Split(role, ",")
	return user, nil
}
