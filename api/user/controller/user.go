/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"net/http"
	"strconv"

	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/response"
	"github.com/TestAndWin/e2e-coverage/user/model"
	"github.com/gin-gonic/gin"
)

const USER_ID = "userId"

// GetUser godoc
// @Summary      Get all user
// @Description  Get all user
// @Tags         user
// @Produce      json
// @Success      200  {object}  model.User
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/users [GET]
func GetUser(c *gin.Context) {
	userStore, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	users, err := userStore.GetUser()
	if err != nil {
		errors.HandleError(c, err)
	} else {
		response.OK(c, users)
	}
}

// CreateUser godoc
// @Summary      Add a new user
// @Description  Takes a user JSON and stores it in DB. Return saved JSON.
// @Tags         user
// @Produce      json
// @Param        user  body      model.User  true  "User JSON"
// @Success      201  {object}  model.User
// @Failure      400  {string}  ErrorResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/users [POST]
func CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Invalid user data", err))
		return
	}

	userStore, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	id, err := userStore.CreateUser(user)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	user.Id = id
	response.Created(c, user)
}

// UpdateUser godoc
// @Summary      Change the role, name and password of a user
// @Description  Takes a user JSON and updates the user.
// @Tags         user
// @Produce      json
// @Param        id    path      int         true  "User ID"
// @Param        user  body      model.User  true  "User JSON"
// @Success      200  {string}  SuccessResponse
// @Failure      400  {string}  ErrorResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/users/{id} [PUT]
func UpdateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Invalid user data", err))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Invalid user ID", err))
		return
	}

	user.Id = id
	userStore, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	err = userStore.UpdateUser(user)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	response.ResponseWithMessage(c, http.StatusOK, "User updated successfully")
}

// ChangePassword godoc
// @Summary      Password Change
// @Description  Takes the NewPassword JSON and updates the password. Only possible for the current user to change his own password.
// @Tags         user
// @Produce      json
// @Param        newPassword  body      model.NewPassword  true  "NewPassword JSON"
// @Success      200  {string}  SuccessResponse
// @Failure      400  {string}  ErrorResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/users/change-pwd [PUT]
func ChangePassword(c *gin.Context) {
	var pwd model.NewPassword

	if err := c.ShouldBindJSON(&pwd); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Invalid password data", err))
		return
	}

	// Get user ID from the context (set by AuthUser middleware)
	id, exists := c.Get(USER_ID)
	if !exists {
		errors.HandleError(c, errors.NewBadRequestError("User ID not found in context", nil))
		return
	}

	// Convert to int64 if needed
	var userId int64
	switch v := id.(type) {
	case int64:
		userId = v
	case float64:
		userId = int64(v)
	case int:
		userId = int64(v)
	default:
		errors.HandleError(c, errors.NewBadRequestError("Invalid user ID type", nil))
		return
	}

	userStore, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	err = userStore.ChangePassword(userId, pwd.Password, pwd.NewPassword)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	response.ResponseWithMessage(c, http.StatusOK, "Password updated successfully")
}

// DeleteUser godoc
// @Summary      Delete the user
// @Description  Delete the user
// @Tags         user
// @Produce      json
// @Param        id    path      int     true  "User ID"
// @Success      204  {string}  SuccessResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/users/{id} [DELETE]
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Invalid user ID", err))
		return
	}

	userStore, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	err = userStore.DeleteUser(id)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// GenerateApiKey godoc
// @Summary      Generate an API Key
// @Description  Generate an API Key
// @Tags         user
// @Produce      json
// @Success      200  {string}  SuccessResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/users/generate-api-key [POST]
func GenerateApiKey(c *gin.Context) {
	userId := c.GetInt64(USER_ID)
	userStore, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	apiKey, err := userStore.GenerateApiKey(userId)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	response.OK(c, gin.H{"key": apiKey})
}
