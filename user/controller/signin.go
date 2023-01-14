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
	"os"
	"strings"
	"time"

	"github.com/TestAndWin/e2e-coverage/config"
	"github.com/TestAndWin/e2e-coverage/user/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = getJwtKey()

func getJwtKey() []byte {
	key, err := config.LoadConfig()
	if err != nil {
		log.Print("Cannot load config ", err)
		os.Exit(0)
	}
	return []byte(key.JWTKey)
}

// Signin godoc
// @Summary      Sign in of a user
// @Description  Sign in and returning a JWT token and a refresh token if user name and password are correct
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
		// Check Login and get role
		roles, err := userStore.Login(s.Email, s.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
			return
		}

		// Create the JWT token
		claims := createClaims(s.Email, roles, time.Now().Add(5*time.Minute))
		token, err := createToken(claims, jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		claims = createClaims(s.Email, roles, time.Now().Add(24*time.Hour))
		refreshToken, err := createToken(claims, jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "roles": roles, "refreshToken": refreshToken})
	}
}

// RefreshToken godoc
// @Summary      Refresh the bearer token
// @Description  Checks if the token is valid and returns the "fresh" the token
// @Tags         user
// @Param        token  body      string  true  "JSON"
// @Produce      json
// @Success      200  {object}  string
// @Router       /api/v1/auth/refresh [POST]
func RefreshToken(c *gin.Context) {
	// Get the token from the post body
	var refresh struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&refresh); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the claims and check them
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(refresh.Token, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Create new token
	newClaims := createClaims(claims.Email, claims.Role, time.Now().Add(5*time.Minute))
	newToken, err := createToken(newClaims, []byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// Create the claims
func createClaims(email string, roles string, t time.Time) *model.Claims {
	claims := &model.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
			Issuer:    "TestAndWin.net",
		},
		Role: roles,
	}
	return claims
}

// Create the token
func createToken(claims *model.Claims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	return signedToken, err
}

// Retrieve token from request header Authorization
func GetToken(c *gin.Context) (*jwt.Token, error) {
	token := c.Request.Header["Authorization"]
	//log.Printf("Auth user: %s", token)
	if token == nil || len(token) < 1 || !strings.Contains(token[0], "Bearer ") {
		return nil, nil
	}
	t := strings.Split(token[0], " ")[1]
	//log.Printf("Token: %s", t)

	// Parse the JWT string
	claims := &model.Claims{}
	return jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })
}

// Reads the Bearer token and checks if the token is valid and if the role saved in the token matches the expected level.
// Level can also be an empty string, then this check is skipped.
func AuthUser(level string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn, err := GetToken(c)
		if tkn == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// Check if token is valid and not expired
		if claims, ok := tkn.Claims.(*model.Claims); ok && tkn.Valid {
			// Has the user the needed role
			if level != "" && strings.Contains(claims.Role, level) {
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
