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

	"github.com/TestAndWin/e2e-coverage/coverage/model"
	"github.com/gin-gonic/gin"
)

// AddFeature godoc
// @Summary      Add a new feature to an area
// @Description  Takes a feature JSON and stores it in DB. Return saved JSON.
// @Tags         feature
// @Produce      json
// @Param        feature  body      model.Feature  true  "Feature JSON"
// @Success      201  {object}  model.Feature
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/features [POST]
func AddFeature(c *gin.Context) {
	var f model.Feature
	if err := c.BindJSON(&f); err != nil {
		handleError(c, err, "Error binding JSON", http.StatusBadRequest)
		return
	}
	id, err := repo.InsertFeature(f)
	if err != nil {
		handleError(c, err, "Error insert feature", http.StatusInternalServerError)
		return
	}
	f.Id = id
	c.JSON(http.StatusCreated, f)
}

// GetAreaFeatures godoc
// @Summary      Get all area features
// @Description  Get all features for the specified area
// @Tags         feature
// @Produce      json
// @Param        id    path      string     true  "Area ID"
// @Success      200  {array}  model.Feature
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/areas/{id}/features [get]
func GetAreaFeatures(c *gin.Context) {
	f, err := repo.GetAllAreaFeatures(c.Param("id"))
	if err != nil {
		handleError(c, err, "Error getting area featues", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, f)
}

// UpdateFeature godoc
// @Summary      Update a feature
// @Description  Takes a feature JSON and feature ID and updates it in DB. Return saved JSON.
// @Tags         feature
// @Param        id       path      int            true  "Feature ID"
// @Param        feature  body      model.Feature  true  "Feature JSON"
// @Produce      json
// @Success      200  {object}  model.Feature
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/features [PUT]
func UpdateFeature(c *gin.Context) {
	var f model.Feature
	if err := c.BindJSON(&f); err != nil {
		handleError(c, err, "Error binding JSON", http.StatusBadRequest)
		return
	}

	f.Id, _ = strconv.ParseInt(c.Param("id"), 0, 64)
	_, err := repo.UpdateFeature(f)
	if err != nil {
		handleError(c, err, "Error updating feature", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, f)
}

// DeleteFeature godoc
// @Summary      Delete the product feature
// @Description  Delete the product feature
// @Tags         feature
// @Produce      json
// @Param        id    path      int     true  "Feature ID"
// @Success      204  {string}  SuccessResponse
// @Failure      500  {string}  ErrorResponse
// @Router       /api/v1/features/{id} [DELETE]
func DeleteFeature(c *gin.Context) {
	_, err := repo.DeleteFeature(c.Param("id"))
	if err != nil {
		handleError(c, err, "Error deleting feature", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"status": "ok"})
}
