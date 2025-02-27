/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/TestAndWin/e2e-coverage/config"
	"github.com/TestAndWin/e2e-coverage/errors"
	"github.com/TestAndWin/e2e-coverage/user/model"
	"github.com/golang-jwt/jwt/v4"
)

// Constants for cookie and context keys
const (
	CookiePayload   = "header.payload"
	CookieSignature = "signature"
	ContextUserID   = "userId"
	ContextUserEmail = "userEmail"
	
	// Token expiry times
	AccessTokenExpiry  = 24 * time.Hour
	RefreshTokenExpiry = 7 * 24 * time.Hour
	
	// Minimum token length for validation
	MinTokenLength = 10
)

// TokenManager handles JWT token creation and validation
type TokenManager struct {
	secretKey     []byte
	refreshSecret []byte
}

// NewTokenManager creates a new token manager with the configured secret keys
func NewTokenManager() (*TokenManager, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to load config: %w", err))
	}
	
	// Generate a separate refresh token secret with a different value
	refreshSecret := make([]byte, 32)
	_, err = rand.Read(refreshSecret)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("failed to generate refresh token secret: %w", err))
	}
	
	return &TokenManager{
		secretKey:     []byte(config.JWTKey),
		refreshSecret: refreshSecret,
	}, nil
}

// CreateAccessToken generates a JWT access token for a user
func (tm *TokenManager) CreateAccessToken(user model.User) (string, error) {
	claims := &model.Claims{
		ID:    user.Id,
		Email: user.Email,
		Role:  user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "TestAndWin.net",
			Subject:   fmt.Sprintf("%d", user.Id),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(tm.secretKey)
	if err != nil {
		return "", errors.NewInternalError(fmt.Errorf("failed to sign access token: %w", err))
	}
	
	return tokenString, nil
}

// CreateRefreshToken generates a refresh token for extending sessions
func (tm *TokenManager) CreateRefreshToken(userID int64) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "TestAndWin.net",
		Subject:   fmt.Sprintf("%d", userID),
		ID:        generateTokenID(), // Add a unique jti (JWT ID) for token identification
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString(tm.refreshSecret)
	if err != nil {
		return "", errors.NewInternalError(fmt.Errorf("failed to sign refresh token: %w", err))
	}
	
	return tokenString, nil
}

// SplitToken splits a JWT token into its header.payload and signature parts
func (tm *TokenManager) SplitToken(token string) (headerPayload, signature string, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", "", errors.NewAppError(
			fmt.Errorf("invalid token format"),
			"Invalid token format",
			"INVALID_TOKEN_FORMAT",
			401,
		)
	}
	
	return parts[0] + "." + parts[1], parts[2], nil
}

// ValidateToken validates a JWT token and returns its claims
func (tm *TokenManager) ValidateToken(tokenString string) (*model.Claims, error) {
	// Basic validation before parsing
	if len(tokenString) < MinTokenLength {
		return nil, errors.NewAppError(
			fmt.Errorf("token too short"),
			"Invalid token",
			"INVALID_TOKEN",
			401,
		)
	}
	
	claims := &model.Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
	})
	
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && validationErr.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, errors.NewAppError(
				err,
				"Token has expired",
				"TOKEN_EXPIRED",
				401,
			)
		}
		return nil, errors.NewAppError(
			err,
			"Invalid token",
			"INVALID_TOKEN",
			401,
		)
	}
	
	if !token.Valid {
		return nil, errors.NewAppError(
			fmt.Errorf("invalid token claims"),
			"Invalid token",
			"INVALID_TOKEN",
			401,
		)
	}
	
	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (tm *TokenManager) ValidateRefreshToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.refreshSecret, nil
	})
	
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && validationErr.Errors&jwt.ValidationErrorExpired != 0 {
			return 0, errors.NewAppError(
				err,
				"Refresh token has expired",
				"REFRESH_TOKEN_EXPIRED",
				401,
			)
		}
		return 0, errors.NewAppError(
			err,
			"Invalid refresh token",
			"INVALID_REFRESH_TOKEN",
			401,
		)
	}
	
	if !token.Valid {
		return 0, errors.NewAppError(
			fmt.Errorf("invalid refresh token"),
			"Invalid refresh token",
			"INVALID_REFRESH_TOKEN",
			401,
		)
	}
	
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errors.NewAppError(
			fmt.Errorf("invalid refresh token claims"),
			"Invalid refresh token",
			"INVALID_REFRESH_TOKEN",
			401,
		)
	}
	
	// Extract user ID from the subject claim
	var userID int64
	_, err = fmt.Sscanf(claims.Subject, "%d", &userID)
	if err != nil {
		return 0, errors.NewAppError(
			fmt.Errorf("invalid subject in refresh token: %w", err),
			"Invalid refresh token",
			"INVALID_REFRESH_TOKEN",
			401,
		)
	}
	
	return userID, nil
}

// GetSecretKey returns the secret key used for signing tokens
// This is needed for token validation
func (tm *TokenManager) GetSecretKey() []byte {
	return tm.secretKey
}

// Helper function to generate a random token ID
func generateTokenID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)
}