/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/TestAndWin/e2e-coverage/user/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// TODO Get this key from the config
var jwtKey = []byte("secret-key")

// TODO Replace with database access later
var users = map[string]string{
	"michael": "passQQQQ",
}

// Signin godoc
// @Summary      Signin of a user
// @Description  Signin and returning a JWT token if user name and password are correct
// @Tags         user
// @Produce      json
// @Param        signin  body      model.Credentials  true  "Credentials JSON"
// @Success      200  {object}  string
// @Router       /api/v1/auth/signin [POST]
func Signin(c *gin.Context) {
	var s model.Credentials
	if err := c.BindJSON(&s); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {

		// TODO Get this from the database, but first hash the password!
		expectedPassword, ok := users[s.Email]
		//roles := model.EDITOR
		roles := fmt.Sprintf("%s,%s,%s", model.ADMIN, model.CONSUMER, model.EDITOR)
		if !ok || expectedPassword != s.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
			return
		}

		// Create the JWT token
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &model.Claims{
			Email: s.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
			Role: roles,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString, "roles": roles})
	}
}

// Retrieve token from request header Authorization
func getToken(c *gin.Context) (*jwt.Token, error) {
	token := c.Request.Header["Authorization"]
	log.Println("Auth user: ", token)
	if token == nil {
		return nil, nil
	}
	t := strings.Split(token[0], " ")[1]

	// Parse the JWT string
	claims := &model.Claims{}
	return jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })
}

// Reads the Bearer token and checks if the token is valid and if the role saved in the token matches the expected level.
func AuthUser(level string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn, err := getToken(c)
		if tkn == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// Check if token is valid and not expired
		if claims, ok := tkn.Claims.(*model.Claims); ok && tkn.Valid {
			if strings.Contains(claims.Role, level) {
				log.Printf("Sign in success: %v %v %v", claims.Email, claims.Role, claims.ExpiresAt)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User has no access to this resource"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RefreshToken godoc
// @Summary      Refresh the bearer token
// @Description  Checks if the token is valid and refreshes the token when expiry is within 60 seconds and returns the new token
// @Tags         user
// @Produce      json
// @Success      200  {object}  string
// @Router       /api/v1/auth/refresh [POST]
func RefreshToken(c *gin.Context) {
	tkn, err := getToken(c)
	if tkn == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		c.Abort()
		return
	}

	if err != nil && err == jwt.ErrSignatureInvalid || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature or token is not valid"})
		c.Abort()
		return
	}

	// Refresh is only possible within 60 seconds after expire
	claims := tkn.Claims.(*model.Claims)
	if time.Until(claims.ExpiresAt.Time) > 60*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh too late"})
		c.Abort()
		return
	}

	// New token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error freshing the token"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
