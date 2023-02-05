/*
Copyright (c) 2022-2023, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"log"
	"net/http"
	"os"

	"github.com/TestAndWin/e2e-coverage/user/model"
	"github.com/TestAndWin/e2e-coverage/user/repository"
	"github.com/gin-gonic/gin"
)

var userStore = initRepository()

// Set-up the db connection and create the db tables if needed
func initRepository() repository.UserStore {
	userStore, err := repository.NewUserStore()
	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err)
		os.Exit(1)
	}

	err = userStore.CreateUsersTable()
	if err != nil {
		log.Fatalf("Error creating tables: %s", err)
		os.Exit(1)
	}
	return *userStore
}

// GetUser godoc
// @Summary      Get all user
// @Description  Get all user
// @Tags         user
// @Produce      json
// @Success      200  {array}  model.User
// @Router       /api/v1/users [GET]
func GetUser(c *gin.Context) {
	user, err := userStore.GetUser()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// CreateUser godoc
// @Summary      Add a new user
// @Description  Takes a user JSON and stores it in DB. Return saved JSON.
// @Tags         user
// @Produce      json
// @Param        user  body      model.User  true  "User JSON"
// @Success      201  {object}  model.User
// @Router       /api/v1/users [POST]
func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := userStore.InsertUser(user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		user.Id = id
		c.JSON(http.StatusCreated, user)
	}
}

// UpdateUser godoc
// @Summary      Change the role, name and password of a user
// @Description  Takes a user JSON and updates the user.
// @Tags         user
// @Produce      json
// @Param        id    path      int         true  "User ID"
// @Param        user  body      model.User  true  "User JSON"
// @Success      200
// @Router       /api/v1/users/{id} [PUT]
func UpdateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userStore.UpdateUser(c.Param("id"), user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// ChangePassword godoc
// @Summary      Password Change
// @Description  Takes the NewPassword JSON and updates the password. Only possible for the current user to change his own password.
// @Tags         user
// @Produce      json
// @Param        id           path      int                true  "User ID"
// @Param        newPassword  body      model.NewPassword  true  "NewPassword JSON"
// @Success      200
// @Router       /api/v1/users/change-pwd/{id} [PUT]
func ChangePassword(c *gin.Context) {
	var pwd model.NewPassword

	if err := c.ShouldBindJSON(&pwd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Token is already checked before. We need the token to get the email
	tkn, _ := GetToken(c)
	if claims, ok := tkn.Claims.(*model.Claims); ok && tkn.Valid {
		err := userStore.UpdatePassword(claims.Email, pwd.Password, pwd.NewPassword)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting token"})
	}
}

// DeleteUser godoc
// @Summary      Delete the user
// @Description  Delete the user
// @Tags         user
// @Produce      json
// @Param        id    path      int     true  "User ID"
// @Success      200
// @Router       /api/v1/users/{id} [DELETE]
func DeleteUser(c *gin.Context) {
	_, err := userStore.DeleteUser(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusNoContent, gin.H{"status": "ok"})
	}
}
