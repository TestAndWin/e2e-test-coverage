/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TestAndWin/e2e-coverage/coverage/model"
	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/response"
	"github.com/gin-gonic/gin"
)

// AddArea godoc
// @Summary      Add a new area to a product
// @Description  Takes an area JSON and stores it in DB. Return saved JSON.
// @Tags         area
// @Produce      json
// @Param        area  body     model.Area  true  "Area JSON"
// @Success      201  {object}  model.Area
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/v1/areas [POST]
func AddArea(c *gin.Context) {
	var a model.Area
	if err := c.BindJSON(&a); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error binding area JSON", err))
		return
	}

	id, err := getRepository().InsertArea(a)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(fmt.Errorf("failed to insert area: %w", err)))
		return
	}

	a.Id = id
	response.Created(c, a)
}

// GetProductAreas godoc
// @Summary      Get all product areas
// @Description  Get all areas for the specified product
// @Tags         area
// @Produce      json
// @Param        id    path    string     true  "Product ID"
// @Success      200  {array}  model.Area
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/v1/products/{id}/areas [get]
func GetProductAreas(c *gin.Context) {
	productID := c.Param("id")
	
	a, err := getRepository().GetAllProductAreas(productID)
	if err != nil {
		errors.HandleError(c, errors.NewAppError(
			err,
			fmt.Sprintf("Error retrieving areas for product %s", productID),
			"AREAS_RETRIEVAL_ERROR",
			http.StatusInternalServerError,
		))
		return
	}
	
	if len(a) == 0 {
		// Return empty list with consistent format
		response.EmptyList(c)
		return
	}
	
	// Return areas with count
	response.ResponseWithDataAndCount(c, http.StatusOK, a, int64(len(a)))
}

// UpdateArea godoc
// @Summary      Update an area
// @Description  Takes an area JSON and the Area ID and updates an area in the DB.
// @Tags         area
// @Produce      json
// @Param        id    path     int         true  "Area ID"
// @Param        area  body     model.Area  true  "Area JSON"
// @Success      200  {object}  model.Area
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/v1/areas [PUT]
func UpdateArea(c *gin.Context) {
	var a model.Area
	if err := c.BindJSON(&a); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error binding area JSON", err))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Invalid area ID", err))
		return
	}
	a.Id = id

	affected, err := getRepository().UpdateArea(a)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(fmt.Errorf("failed to update area %d: %w", id, err)))
		return
	}
	
	// Check if the area was found
	if affected == 0 {
		errors.HandleError(c, errors.NewNotFoundError(fmt.Sprintf("Area with ID %d", id)))
		return
	}
	
	response.ResponseWithDataAndMessage(c, http.StatusOK, a, "Area updated successfully")
}

// DeleteArea godoc
// @Summary      Delete the product area
// @Description  Delete the product area
// @Tags         area
// @Produce      json
// @Param        id    path      int     true  "Area ID"
// @Success      204  {string}  SuccessResponse
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      404  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/v1/areas/{id} [DELETE]
func DeleteArea(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Invalid area ID", err))
		return
	}

	affected, err := getRepository().DeleteArea(id)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(fmt.Errorf("failed to delete area %d: %w", id, err)))
		return
	}
	
	// Check if the area was found
	if affected == 0 {
		errors.HandleError(c, errors.NewNotFoundError(fmt.Sprintf("Area with ID %d", id)))
		return
	}

	response.NoContent(c)
}
