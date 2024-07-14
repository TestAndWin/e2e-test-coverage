package controller

import (
	"log"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error, message string, status int) {
	log.Printf("%s: %v", message, err)
	c.JSON(status, gin.H{
		"error":   message,
		"details": err.Error(),
		"status":  status,
	})
}
