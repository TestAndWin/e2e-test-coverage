/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"net/http"

	"github.com/TestAndWin/e2e-coverage/coverage/model"
	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/response"
	"github.com/TestAndWin/e2e-coverage/user/controller"
	"github.com/gin-gonic/gin"
)

// AddTest godoc
// @Summary      Add a new expl test
// @Description  Takes a exploratory test JSON and stores it in DB. Return saved JSON.
// @Tags         expl-test
// @Produce      json
// @Param        expl-test  body      model.ExplTest  true  "Expl Test JSON"
// @Success      201  {object}  model.ExplTest
// @Failure      400  {string}  ErrorResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/expl-tests [POST]
func AddExplTest(c *gin.Context) {
	var et model.ExplTest
	if err := c.BindJSON(&et); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error binding JSON", err))
		return
	}

	et.Tester = c.GetInt64(controller.USER_ID)
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	id, err := repo.InsertExplTest(et)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}

	et.Id = id
	response.Created(c, et)
}

// DeleteTest godoc
// @Summary      Delete an expl test
// @Description  Delete an expl test
// @Tags         expl-test
// @Produce      json
// @Param        id    path      int     true  "Test ID"
// @Success      204  {string} SuccesResponse
// @Failure      500  {string} ErrorResponse
// @Router       /api/v1/expl-tests/{id} [DELETE]
func DeleteExplTest(c *gin.Context) {
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	_, err = repo.DeleteExplTest(c.Param("id"))
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// GetExplTestsForArea godoc
// @Summary      Get all exploratory tests.
// @Description  Get all exploratory tests for the specified area for the last 28 days
// @Tags         expl-test
// @Produce      json
// @Param        areaid    path      int     true  "Area ID"
// @Success      200  {array}  model.ExplTest
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/expl-tests/area/{areaid} [POST]
func GetExplTestsForArea(c *gin.Context) {
	repo, err := getRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	et, err := repo.GetExplTests(c.Param("areaid"))
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	response.OK(c, et)
}
