/*
Copyright (c) 2022-2024, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DeleteTests godoc
// @Summary      Delete all tests for the specified component, suite and file-name
// @Description  Delete all tests for the specified component, suite and file-name
// @Tags         test
// @Produce      json
// @Param        component      query      string     true  "Component name"
// @Param        suite          query      string     true  "Suite name"
// @Param        file-name      query      string     true  "File name"
// @Success      204 {object}  SucccessResponse
// @Success      500 {object}  ErrorResponse
// @Router       /api/v1/tests/{id} [DELETE]
func DeleteTests(c *gin.Context) {
	suite := c.Query("suite")
	component := c.Query("component")
	file := strings.Replace(c.Query("file-name"), "\\\\", "\\", -1)
	_, err := repo.DeleteTest(component, suite, file)
	if err != nil {
		handleError(c, err, "Error delete tests", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"status": "ok"})
}

// GetAllTestForSuiteFile godoc
// @Summary      Get all tests for the specified suite and filename.
// @Description  Get all tests for the specified suite and filename.
// @Tags         test
// @Produce      json
// @Param        component      query      string     true  "Component name"
// @Param        suite          query      string     true  "Suite name"
// @Param        file-name      query      string     true  "File name"
// @Success      200 {array}  model.Test
// @Success      500 {object} ErrorResponse
// @Router       /api/v1/tests [GET]
func GetAllTestForSuiteFile(c *gin.Context) {
	suite := c.Query("suite")
	component := c.Query("component")
	file := strings.Replace(c.Query("file-name"), "\\\\", "\\", -1)
	tests, err := repo.GetAllTestForSuiteFile(component, suite, file)
	if err != nil {
		handleError(c, err, "Error getting tests", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, tests)
}
