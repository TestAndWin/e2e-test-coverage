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

	"github.com/TestAndWin/e2e-coverage/pkg/model"
	"github.com/gin-gonic/gin"
)

// AddTest godoc
// @Summary      Add a new expl test
// @Description  Takes a exploratory test JSON and stores it in DB. Return saved JSON.
// @Tags         expl-test
// @Produce      json
// @Param        expl-test  body      model.ExplTest  true  "Expl Test JSON"
// @Success      201  {object}  model.ExplTest
// @Router       /api/v1/expl-tests [POST]
func AddExplTest(c *gin.Context) {
	var et model.ExplTest
	if err := c.BindJSON(&et); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		id, err := repo.ExecuteSql(model.INSERT_EXPL_TEST, et.AreaId, et.Summary, et.Rating, et.TestRun)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": err})
		} else {
			et.Id = id
			c.JSON(http.StatusCreated, et)
		}
	}
}

// DeleteTest godoc
// @Summary      Deletes the expl test
// @Description  Deletes the expl test
// @Tags         expl-test
// @Produce      json
// @Param        id    path      int     true  "Test ID"
// @Success      200
// @Router       /api/v1/expl-tests/{id} [DELETE]
func DeleteExplTest(c *gin.Context) {
	_, err := repo.ExecuteSql(model.DELETE_TEST, c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusNoContent, gin.H{"status": "ok"})
	}
}

// GetExplTestsForArea godoc
// @Summary      Get all exploratory tests.
// @Description  Get all exploratory tests for the specified area for the last 28 days
// @Tags         expl-test
// @Produce      json
// @Param        areaid    path      int     true  "Area ID"
// @Success      200  {array}  model.ExplTest
// @Router       /api/v1/expl-tests/area/{areaid} [POST]
func GetExplTestsForArea(c *gin.Context) {
	et, err := repo.GetExplTests(c.Param("areaid"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, et)
	}
}
