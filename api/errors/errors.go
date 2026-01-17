/*
Copyright (c) 2022-2026, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package errors

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// AppError represents a structured application error
type AppError struct {
	// Original error
	Err error
	// User-friendly message
	Message string
	// Optional error code
	Code string
	// HTTP status code
	StatusCode int
	// Stack trace
	Stack string
}

// Error satisfies the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// ErrorResponse is the structured JSON response for errors
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Status  int    `json:"status"`
	Details string `json:"details,omitempty"`
}

// Common error types
var (
	ErrNotFound     = NewAppError(errors.New("resource not found"), "Not Found", "NOT_FOUND", http.StatusNotFound)
	ErrBadRequest   = NewAppError(errors.New("invalid request"), "Bad Request", "BAD_REQUEST", http.StatusBadRequest)
	ErrInternal     = NewAppError(errors.New("internal server error"), "Internal Server Error", "INTERNAL_ERROR", http.StatusInternalServerError)
	ErrForbidden    = NewAppError(errors.New("forbidden"), "Forbidden", "FORBIDDEN", http.StatusForbidden)
	ErrUnauthorized = NewAppError(errors.New("unauthorized"), "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized)
)

// NewAppError creates a new application error
func NewAppError(err error, message, code string, statusCode int) *AppError {
	// Capture stack trace
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	// Format stack trace
	var stackTrace string
	for {
		frame, more := frames.Next()
		stackTrace += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}

	return &AppError{
		Err:        err,
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
		Stack:      stackTrace,
	}
}

// NewNotFoundError creates a not found error with custom message
func NewNotFoundError(resource string) *AppError {
	return NewAppError(
		fmt.Errorf("%s not found", resource),
		fmt.Sprintf("%s not found", resource),
		"NOT_FOUND",
		http.StatusNotFound,
	)
}

// NewBadRequestError creates a bad request error with custom message
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(
		err,
		message,
		"BAD_REQUEST",
		http.StatusBadRequest,
	)
}

// NewInternalError creates an internal server error with optional error
func NewInternalError(err error) *AppError {
	return NewAppError(
		err,
		"Internal Server Error",
		"INTERNAL_ERROR",
		http.StatusInternalServerError,
	)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(
		errors.New("unauthorized"),
		message,
		"UNAUTHORIZED",
		http.StatusUnauthorized,
	)
}

// HandleError handles an error and responds to the client
func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		// Convert standard error to AppError
		appErr = NewInternalError(err)
	}

	// Log the error with stack trace in development/staging
	log.Printf("[ERROR] %s: %v\nStack:\n%s", appErr.Message, appErr.Err, appErr.Stack)

	// Send structured response to client (without stack trace)
	c.JSON(appErr.StatusCode, ErrorResponse{
		Success: false,
		Message: appErr.Message,
		Code:    appErr.Code,
		Status:  appErr.StatusCode,
		Details: appErr.Err.Error(),
	})
}
