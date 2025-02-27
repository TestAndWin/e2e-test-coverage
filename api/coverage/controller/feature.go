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
	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/response"
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
		errors.HandleError(c, errors.NewBadRequestError("Error binding JSON", err))
		return
	}
	
	repo := getRepository()
	id, err := repo.InsertFeature(f)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	f.Id = id
	response.Created(c, f)
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
	areaID := c.Param("id")
	
	repo := getRepository()
	features, err := repo.GetAllAreaFeatures(areaID)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	
	// Use ResponseWithDataAndCount to ensure consistent format
	// even for empty arrays
	if len(features) == 0 {
		// Return empty list with consistent format
		response.EmptyList(c)
	} else {
		// Return features with count
		response.ResponseWithDataAndCount(c, http.StatusOK, features, int64(len(features)))
	}
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
		errors.HandleError(c, errors.NewBadRequestError("Error binding JSON", err))
		return
	}

	repo := getRepository()
	f.Id, _ = strconv.ParseInt(c.Param("id"), 0, 64)
	_, err := repo.UpdateFeature(f)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	response.OK(c, f)
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
	repo := getRepository()
	_, err := repo.DeleteFeature(c.Param("id"))
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	c.Status(http.StatusNoContent)
}
