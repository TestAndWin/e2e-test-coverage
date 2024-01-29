/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"log"
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
// @Router       /api/v1/features [POST]
func AddFeature(c *gin.Context) {
	var f model.Feature
	if err := c.BindJSON(&f); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		id, err := repo.InsertFeature(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			f.Id = id
			c.JSON(http.StatusCreated, f)
		}
	}
}

// GetAreaFeatures godoc
// @Summary      Get all area features
// @Description  Get all features for the specified area
// @Tags         feature
// @Produce      json
// @Param        id    path      string     true  "Area ID"
// @Success      200  {array}  model.Feature
// @Router       /api/v1/areas/{id}/features [get]
func GetAreaFeatures(c *gin.Context) {
	f, err := repo.GetAllAreaFeatures(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, f)
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
// @Router       /api/v1/features [PUT]
func UpdateFeature(c *gin.Context) {
	var f model.Feature
	if err := c.BindJSON(&f); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		f.Id, _ = strconv.ParseInt(c.Param("id"), 0, 64)
		_, err := repo.UpdateFeature(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			c.JSON(http.StatusOK, f)
		}
	}
}

// DeleteFeature godoc
// @Summary      Delete the product feature
// @Description  Delete the product feature
// @Tags         feature
// @Produce      json
// @Param        id    path      int     true  "Feature ID"
// @Success      200
// @Router       /api/v1/features/{id} [DELETE]
func DeleteFeature(c *gin.Context) {
	_, err := repo.DeleteFeature(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
