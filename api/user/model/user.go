/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

// User represents a user in the system
type User struct {
	Id       int64    `json:"id"       db:"id"`
	Email    string   `json:"email"    db:"email"`
	Password string   `json:"password" db:"password"`
	Roles    []string `json:"roles"    db:"roles"`
}

// Struct to update the password
type NewPassword struct {
	Password    string `json:"password"`
	NewPassword string `json:"new-password"`
}

// Can view the coverage and report new exploratory tests
const TESTER = "Tester"

// Can create new products, products areas, ...
const MAINTAINER = "Maintainer"

// Can create new user and edit them
const ADMIN = "Admin"
