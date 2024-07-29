/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"net/http"
	"strconv"

	"github.com/TestAndWin/e2e-coverage/coverage/model"
	"github.com/gin-gonic/gin"
)

// AddArea godoc
// @Summary      Add a new area to a product
// @Description  Takes an area JSON and stores it in DB. Return saved JSON.
// @Tags         area
// @Produce      json
// @Param        area  body     model.Area  true  "Area JSON"
// @Success      201  {object}  model.Area
// @Failure      400  {string}  Invalid JSON
// @Failure      500  {string}  Internal Error
// @Router       /api/v1/areas [POST]
func AddArea(c *gin.Context) {
	var a model.Area
	if err := c.BindJSON(&a); err != nil {
		handleError(c, err, "Error binding JSON", http.StatusBadRequest)
		return
	}

	id, err := repo.InsertArea(a)
	if err != nil {
		handleError(c, err, "Error insert area", http.StatusInternalServerError)
		return
	}

	a.Id = id
	c.JSON(http.StatusCreated, a)
}

// GetProductAreas godoc
// @Summary      Get all product areas
// @Description  Get all areas for the specified product
// @Tags         area
// @Produce      json
// @Param        id    path    string     true  "Product ID"
// @Success      200  {array}  model.Area
// @Failure      500  {string} Internal Error
// @Router       /api/v1/products/{id}/areas [get]
func GetProductAreas(c *gin.Context) {
	a, err := repo.GetAllProductAreas(c.Param("id"))
	if err != nil {
		handleError(c, err, "Error getting areas", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, a)
}

// UpdateArea godoc
// @Summary      Update an area
// @Description  Takes an area JSON and the Area ID and updates an area in the DB.
// @Tags         area
// @Produce      json
// @Param        id    path     int         true  "Area ID"
// @Param        area  body     model.Area  true  "Area JSON"
// @Success      200  {object}  model.Area
// @Failure      400  {string}  Invalid JSON or ID
// @Failure      500  {string}  Internal Error
// @Router       /api/v1/areas [PUT]
func UpdateArea(c *gin.Context) {
	var a model.Area
	if err := c.BindJSON(&a); err != nil {
		handleError(c, err, "Error binding JSON", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		handleError(c, err, "Invalid Area ID", http.StatusBadRequest)
		return
	}
	a.Id = id

	_, err = repo.UpdateArea(a)
	if err != nil {
		handleError(c, err, "Error updating area", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, a)
}

// DeleteArea godoc
// @Summary      Delete the product area
// @Description  Delete the product area
// @Tags         area
// @Produce      json
// @Param        id    path      int     true  "Area ID"
// @Success      204  {string}  SuccessResponse
// @Failure      400  {string}  ErrorResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/areas/{id} [DELETE]
func DeleteArea(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		handleError(c, err, "Invalid Area ID", http.StatusBadRequest)
		return
	}

	_, err = repo.DeleteArea(id)
	if err != nil {
		handleError(c, err, "Error deleting area", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"status": "ok"})
}
