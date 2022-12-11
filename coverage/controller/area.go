/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
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

// AddArea godoc
// @Summary      Add a new area to a product
// @Description  Takes an area JSON and stores it in DB. Return saved JSON.
// @Tags         area
// @Produce      json
// @Param        area  body      model.Area  true  "Area JSON"
// @Success      201  {object}  model.Area
// @Router       /api/v1/areas [POST]
func AddArea(c *gin.Context) {
	var a model.Area
	if err := c.BindJSON(&a); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		id, err := repo.InsertArea(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			a.Id = id
			c.JSON(http.StatusCreated, a)
		}
	}
}

// GetProductAreas godoc
// @Summary      Get all product areas
// @Description  Get all areas for the specified product
// @Tags         area
// @Produce      json
// @Param        id    path      string     true  "Product ID"
// @Success      200  {array}  model.Area
// @Router       /api/v1/products/{id}/areas [get]
func GetProductAreas(c *gin.Context) {
	a, err := repo.GetAllProductAreas(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, a)
	}
}

// UpdateArea godoc
// @Summary      Updates an area
// @Description  Takes an area JSON and the Area ID and updates an area in the DB.
// @Tags         area
// @Produce      json
// @Param        id    path      int         true  "Area ID"
// @Param        area  body      model.Area  true  "Area JSON"
// @Success      200  {object}  model.Area
// @Router       /api/v1/areas [PUT]
func UpdateArea(c *gin.Context) {
	var a model.Area
	if err := c.BindJSON(&a); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		a.Id, _ = strconv.ParseInt(c.Param("id"), 0, 64)
		_, err := repo.UpdateArea(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
		} else {
			c.JSON(http.StatusOK, a)
		}
	}
}

// DeleteArea godoc
// @Summary      Deletes the product area
// @Description  Deletes the product area
// @Tags         area
// @Produce      json
// @Param        id    path      int     true  "Area ID"
// @Success      200
// @Router       /api/v1/areas/{id} [DELETE]
func DeleteArea(c *gin.Context) {
	_, err := repo.DeleteArea(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
