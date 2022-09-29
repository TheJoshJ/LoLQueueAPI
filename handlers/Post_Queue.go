package handlers

import "github.com/gin-gonic/gin"

func Queue(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Data Received!",
	})
}
