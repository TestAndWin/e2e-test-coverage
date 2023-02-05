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
	"strconv"
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
	role VARCHAR(50),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

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
	return s.executeSql("INSERT INTO users(email, password, role) VALUES (?,?,?)", user.Email, hashedPassword, strings.Join(user.Roles, ","))
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
func (s *UserStore) UpdatePassword(email string, password string, newPassword string) error {
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
	_, err = s.executeSql("UPDATE users SET password = ? WHERE email = ?", hashedPassword, email)
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

// Returns all user
func (s UserStore) GetUser() ([]model.User, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := s.db.PrepareContext(ctx, "SELECT id, email, role FROM users;")
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
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
func (s UserStore) DeleteUser(id string) (int64, error) {
	return s.executeSql("DELETE FROM users WHERE id = ?", id)
}

// Creates a new API Key
func (s UserStore) GenerateApiKey(userId int64) (string, error) {
	// TODO Store API Key in users Table - extend table with new column
	return strconv.FormatInt(int64(userId), 10), nil
}
