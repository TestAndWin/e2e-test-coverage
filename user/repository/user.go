/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package repository

import (
	"context"
	"database/sql"
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
	role VARCHAR(5),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

const insertUserStmt = "INSERT INTO users(email, password, role) VALUES (?,?,?)"
const updatePasswordStmt = "UPDATE users SET password = ? WHERE id = ? and email = ?"
const loginStmt = "SELECT password, role FROM users WHERE email = ?"

// UserStore is a store for user data that uses a MySQL database
type UserStore struct {
	db *sql.DB
}

// NewUserStore creates a new UserStore that uses a MySQL database
func NewUserStore() (*UserStore, error) {
	db, err := db.OpenDbConnection("user")
	if err != nil {
		return nil, err
	}
	return &UserStore{db: db}, nil
}

// Create the users table if not existing and create an admin user
func (s *UserStore) CreateUsersTable() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

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
		_, err = s.InsertUser(user)
	}
	return err
}

// Inserts a new user. The password will be stored encrypted. It returns the id of the new user.
func (s *UserStore) InsertUser(user model.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// We keep it simple for now
	return s.executeSql(insertUserStmt, user.Email, hashedPassword, strings.Join(user.Roles, ","))
}

// Updates the user
func (s *UserStore) UpdateUser(id string, user model.User) error {
	if len(user.Email) > 0 {
		_, err := s.executeSql("UPDATE users SET email = ? WHERE id = ?", user.Email, id)
		if err != nil {
			return err
		}
	}

	_, err := s.executeSql("UPDATE users SET role = ? WHERE id = ?", strings.Join(user.Roles, ","), id)
	if err != nil {
		return err
	}

	if len(user.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = s.executeSql("UPDATE users SET password = ? WHERE id = ?", hashedPassword, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// Updates the password. As first step the login with the old password is validated.
func (s *UserStore) UpdatePassword(id string, email string, password string, newPassword string) error {
	// First check the old password
	_, err := s.Login(email, password)
	if err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password
	_, err = s.executeSql(updatePasswordStmt, hashedPassword, id, email)
	return err
}

// Checks if the user can login and if yes, returns the roles, otherwise an error
func (s *UserStore) Login(email string, password string) (string, error) {
	var role string
	var hashedPassword []byte
	err := s.db.QueryRow(loginStmt, email).Scan(&hashedPassword, &role)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return "", err
	}

	return role, nil
}

// Inserts/Deletes a row using the specified statement and params
func (s *UserStore) executeSql(statement string, params ...any) (int64, error) {
	log.Printf("statement %s, params: %s", statement, params)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := s.db.PrepareContext(ctx, statement)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, params...)
	if err != nil {
		log.Printf("Error %s when executing statement", err)
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
