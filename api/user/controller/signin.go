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

const COOKIE_PAYLOAD = "header.payload"
const COOKIE_SIGNATURE = "signature"
const USER_ID = "userId"

var jwtKey = getJwtKey()

func getJwtKey() []byte {
	key, err := config.LoadConfig()
	if err != nil {
		log.Print("Cannot load config ", err)
		os.Exit(0)
	}
	return []byte(key.JWTKey)
}

// Login godoc
// @Summary      Log in of a user
// @Description  Log in and returning a JWT token and a refresh token if user name and password are correct
// @Tags         user
// @Produce      json
// @Param        login  body      model.Credentials  true  "Credentials JSON"
// @Success      200  {object}  string
// @Router       /api/v1/auth/login [POST]
func Login(c *gin.Context) {
	var s model.Credentials
	if err := c.BindJSON(&s); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// Check Login and get role
		user, err := userStore.Login(s.Email, s.Password)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"status": "Unauthorized"})
			return
		}

		// Create the JWT token
		claims := createClaims(user, time.Now().Add(24*time.Hour))
		token, err := createToken(claims, jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Split JWT Token Approach
		t := strings.Split(token, ".")
		c.SetCookie(COOKIE_PAYLOAD, t[0]+"."+t[1], 1800, "/", "", true, false)

		// signature - Session live
		cookie := &http.Cookie{
			Name:     COOKIE_SIGNATURE,
			Value:    t[2],
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(c.Writer, cookie)
		c.JSON(http.StatusOK, gin.H{"roles": strings.Join(user.Roles, ",")})
	}
}

// Create the claims
func createClaims(user model.User, t time.Time) *model.Claims {
	claims := &model.Claims{
		ID:    user.Id,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
			Issuer:    "TestAndWin.net",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: user.Roles,
	}
	return claims
}

// Create the token
func createToken(claims *model.Claims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	return signedToken, err
}

// Retrieve token from cookies
func GetToken(c *gin.Context) (*jwt.Token, error) {
	// Get cookie values
	p, err := c.Cookie(COOKIE_PAYLOAD)
	if err != nil {
		c.String(http.StatusUnauthorized, "Cookie not set")
		return nil, nil
	}
	s, err := c.Cookie(COOKIE_SIGNATURE)
	if err != nil {
		c.String(http.StatusUnauthorized, "Cookie not set")
		return nil, nil
	}
	// And get the complete token
	t := p + "." + s

	// Parse the JWT string
	claims := &model.Claims{}
	return jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil })
}

// Reads the Bearer token and checks if the token is valid and if the role saved in the token matches the expected level.
// Level can also be an empty string, then this check is skipped.
// The header.payload cookie is updated with a new expire date (+30m)
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
			// Has the user the needed role?
			if level == "" || strings.Contains(strings.Join(claims.Role, ","), level) {
				log.Printf("Log in success: %v %v %v", claims.Email, claims.Role, claims.ExpiresAt)
				// Set the user id in case it is needed later
				c.Set(USER_ID, claims.ID)
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "User has no access to this resource"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Update the expire date
		p, _ := c.Cookie(COOKIE_PAYLOAD)
		c.SetCookie(COOKIE_PAYLOAD, p, 1800, "/", "", true, false)
		c.Next()
	}
}

// Reads HTTP Header
func AuthApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("apiKey")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "apiKey is missing"})
			c.Abort()
			return
		}
		userId, err := userStore.GetUserIdForApiKey(apiKey)
		if userId < 1 || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong apiKey"})
			c.Abort()
			return
		}
		log.Printf("Request with API-Key for user %d", userId)
		c.Next()
	}
}
