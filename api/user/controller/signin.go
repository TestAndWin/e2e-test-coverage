/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TestAndWin/e2e-coverage/auth"
	"github.com/TestAndWin/e2e-coverage/dependency"
	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/logger"
	"github.com/TestAndWin/e2e-coverage/response"
	"github.com/TestAndWin/e2e-coverage/user/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// getTokenManager returns the token manager from the dependency container
func getTokenManager() *auth.TokenManager {
	container := dependency.GetContainer()
	manager, err := container.GetTokenManager()
	if err != nil {
		logger.Errorf("Failed to get token manager: %v", err)
		return nil
	}
	return manager
}

// Login godoc
// @Summary      Log in of a user
// @Description  Log in and returning a JWT token and a refresh token if user name and password are correct
// @Tags         user
// @Produce      json
// @Param        login  body      model.Credentials  true  "Credentials JSON"
// @Success      200  {object}  response.StandardResponse
// @Failure      400  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      403  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/v1/auth/login [POST]
func Login(c *gin.Context) {
	var credentials model.Credentials
	if err := c.BindJSON(&credentials); err != nil {
		errors.HandleError(c, errors.NewBadRequestError("Error binding JSON", err))
		return
	}

	// Validate input
	if credentials.Email == "" || credentials.Password == "" {
		errors.HandleError(c, errors.NewBadRequestError("Email and password are required", nil))
		return
	}

	// Check Login and get user with roles
	repo, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	user, err := repo.Login(credentials.Email, credentials.Password)
	if err != nil {
		errors.HandleError(c, errors.NewAppError(
			err,
			"Login failed",
			"AUTH_FAILED",
			http.StatusForbidden,
		))
		return
	}

	// Create the access and refresh tokens
	tm := getTokenManager()

	// Create access token
	accessToken, err := tm.CreateAccessToken(user)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(fmt.Errorf("error creating access token: %w", err)))
		return
	}

	// Create refresh token
	refreshToken, err := tm.CreateRefreshToken(user.Id)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(fmt.Errorf("error creating refresh token: %w", err)))
		return
	}

	// Split access token for security
	headerPayload, signature, err := tm.SplitToken(accessToken)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// Set cookies for access token parts
	c.SetCookie(
		auth.CookiePayload,
		headerPayload,
		int(auth.AccessTokenExpiry.Seconds()),
		"/",
		"",
		true,  // Secure
		false, // Not HttpOnly to allow JavaScript access
	)

	// Signature cookie with HttpOnly flag for security
	cookie := &http.Cookie{
		Name:     auth.CookieSignature,
		Value:    signature,
		Path:     "/",
		MaxAge:   int(auth.AccessTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)

	// Store refresh token in a secure HttpOnly cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh",
		MaxAge:   int(auth.RefreshTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, refreshCookie)

	// For debugging
	logger.Debugf("User login successful: %s, roles: %v", user.Email, user.Roles)

	// Return success response with user info
	response.ResponseWithDataAndMessage(c, http.StatusOK,
		gin.H{
			"userId": user.Id,
			"email":  user.Email,
			"roles":  strings.Join(user.Roles, ","),
		},
		"Login successful",
	)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Use a refresh token to get a new access token
// @Tags         user
// @Produce      json
// @Success      200  {object}  response.StandardResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /api/v1/auth/refresh [POST]
func RefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		errors.HandleError(c, errors.NewUnauthorizedError("Refresh token not found"))
		return
	}

	tm := getTokenManager()

	// Validate refresh token and get user ID
	userID, err := tm.ValidateRefreshToken(refreshToken)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// Get user from database to ensure it still exists and has correct roles
	repo, err := getUserRepository()
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(err))
		return
	}
	user, err := repo.GetUserById(userID)
	if err != nil {
		errors.HandleError(c, errors.NewAppError(
			err,
			"User not found",
			"USER_NOT_FOUND",
			http.StatusUnauthorized,
		))
		return
	}

	// Create new access token
	accessToken, err := tm.CreateAccessToken(user)
	if err != nil {
		errors.HandleError(c, errors.NewInternalError(fmt.Errorf("error creating access token: %w", err)))
		return
	}

	// Split the token and set cookies
	headerPayload, signature, err := tm.SplitToken(accessToken)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// Set cookies for the new access token
	c.SetCookie(
		auth.CookiePayload,
		headerPayload,
		int(auth.AccessTokenExpiry.Seconds()),
		"/",
		"",
		true,
		false,
	)

	cookie := &http.Cookie{
		Name:     auth.CookieSignature,
		Value:    signature,
		Path:     "/",
		MaxAge:   int(auth.AccessTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)

	// Return success response
	response.ResponseWithMessage(c, http.StatusOK, "Token refreshed successfully")
}

// Logout godoc
// @Summary      Log out a user
// @Description  Clear auth cookies to log out the user
// @Tags         user
// @Produce      json
// @Success      200  {object}  response.StandardResponse
// @Router       /api/v1/auth/logout [POST]
func Logout(c *gin.Context) {
	// Clear all authentication cookies
	c.SetCookie(auth.CookiePayload, "", -1, "/", "", true, false)
	c.SetCookie(auth.CookieSignature, "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/api/v1/auth/refresh", "", true, true)

	response.ResponseWithMessage(c, http.StatusOK, "Logged out successfully")
}

// AuthUser middleware checks if the user is authenticated and has the required role
// Level can be an empty string, in which case the role check is skipped
func AuthUser(level string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tm := getTokenManager()

		// Get token parts from cookies
		headerPayload, err := c.Cookie(auth.CookiePayload)
		if err != nil {
			errors.HandleError(c, errors.NewUnauthorizedError("Authentication required"))
			c.Abort()
			return
		}

		signature, err := c.Cookie(auth.CookieSignature)
		if err != nil {
			errors.HandleError(c, errors.NewUnauthorizedError("Authentication required"))
			c.Abort()
			return
		}

		// Combine token parts
		tokenString := headerPayload + "." + signature

		// Validate token
		claims, err := tm.ValidateToken(tokenString)
		if err != nil {
			errors.HandleError(c, err)
			c.Abort()
			return
		}

		// Check role if required
		if level != "" && !hasRole(claims.Role, level) {
			errors.HandleError(c, errors.NewAppError(
				fmt.Errorf("required role: %s", level),
				"Insufficient permissions",
				"ACCESS_DENIED",
				http.StatusForbidden,
			))
			c.Abort()
			return
		}

		// Set user info in context for later use
		c.Set(auth.ContextUserID, claims.ID)
		c.Set(auth.ContextUserEmail, claims.Email)

		// Continue with request
		c.Next()
	}
}

// Helper function to check if user has a specific role
func hasRole(roles []string, requiredRole string) bool {
	for _, role := range roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

// GetToken reconstructs and returns the JWT token from cookies
func GetToken(c *gin.Context) (*jwt.Token, error) {
	// Get token parts from cookies
	headerPayload, err := c.Cookie(auth.CookiePayload)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Authentication required")
	}

	signature, err := c.Cookie(auth.CookieSignature)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Authentication required")
	}

	// Combine token parts
	tokenString := headerPayload + "." + signature

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the token's secret key from the token manager
		tm := getTokenManager()
		secretKey := tm.GetSecretKey()
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// AuthApi middleware for API key authentication
func AuthApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("apiKey")
		if apiKey == "" {
			errors.HandleError(c, errors.NewUnauthorizedError("API key is missing"))
			c.Abort()
			return
		}

		repo, err := getUserRepository()
		if err != nil {
			errors.HandleError(c, errors.NewInternalError(err))
			c.Abort()
			return
		}
		userId, err := repo.GetUserIdForApiKey(apiKey)
		if userId < 1 || err != nil {
			errors.HandleError(c, errors.NewAppError(
				err,
				"Invalid API key",
				"INVALID_API_KEY",
				http.StatusUnauthorized,
			))
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set(auth.ContextUserID, userId)
		logger.Debugf("Request with API-Key for user %d", userId)
		c.Next()
	}
}
