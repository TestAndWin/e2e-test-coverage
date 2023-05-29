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
	"strings"

	"github.com/gin-gonic/gin"
)

// DeleteTest godoc
// @Summary      Delete the test
// @Description  Delete the test
// @Tags         test
// @Produce      json
// @Param        id    path      int     true  "Test ID"
// @Success      204
// @Router       /api/v1/tests/{id} [DELETE]
func DeleteTest(c *gin.Context) {
	_, err := repo.DeleteTest(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusNoContent, gin.H{"status": "ok"})
	}
}

// GetAllTestForSuiteFile godoc
// @Summary      Get all tests for the specified suite and filename.
// @Description  Get all tests for the specified suite and filename.
// @Tags         test
// @Produce      json
// @Param        suite          query      string     true  "Suite name"
// @Param        file-name    query      string     true  "File name"
// @Success      200 {array}  model.Test
// @Router       /api/v1/tests [GET]
func GetAllTestForSuiteFile(c *gin.Context) {
	suite := c.Query("suite")
	file := strings.Replace(c.Query("file-name"), "\\\\", "\\", -1)
	tests, err := repo.GetAllTestForSuiteFile(suite, file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
	} else {
		c.JSON(http.StatusOK, tests)
	}
}
