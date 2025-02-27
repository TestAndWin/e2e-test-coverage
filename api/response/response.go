/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse represents a standardized API response format
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Count   int64       `json:"count,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ResponseWithData returns a standardized successful response with data
func ResponseWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Data:    data,
	})
}

// ResponseWithMessage returns a standardized successful response with message
func ResponseWithMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Message: message,
	})
}

// ResponseWithDataAndMessage returns a standardized successful response with data and message
func ResponseWithDataAndMessage(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ResponseWithDataAndCount returns a standardized successful response with data and count
func ResponseWithDataAndCount(c *gin.Context, statusCode int, data interface{}, count int64) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Data:    data,
		Count:   count,
	})
}

// ResponseWithDataAndMeta returns a standardized successful response with data and metadata
func ResponseWithDataAndMeta(c *gin.Context, statusCode int, data interface{}, meta interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Common HTTP status helpers

// OK sends a 200 OK response with data
func OK(c *gin.Context, data interface{}) {
	ResponseWithData(c, http.StatusOK, data)
}

// Created sends a 201 Created response with data
func Created(c *gin.Context, data interface{}) {
	ResponseWithData(c, http.StatusCreated, data)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// EmptyList returns an empty array with 200 status
func EmptyList(c *gin.Context) {
	ResponseWithDataAndCount(c, http.StatusOK, []interface{}{}, 0)
}

// PaginatedList returns a paginated list with metadata
func PaginatedList(c *gin.Context, data interface{}, count int64, page, pageSize int) {
	meta := map[string]interface{}{
		"page":       page,
		"pageSize":   pageSize,
		"totalItems": count,
		"totalPages": (count + int64(pageSize) - 1) / int64(pageSize),
	}
	ResponseWithDataAndMeta(c, http.StatusOK, data, meta)
}