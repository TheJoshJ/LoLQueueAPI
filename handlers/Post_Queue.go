package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/models"
)

func Queue(c *gin.Context) {
	log.Println(&c)
	var command models.Queue
	err := c.BindJSON(&command)
	if err != nil {
		c.JSON(406, gin.H{
			"message": err,
		})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "Data Received!",
		})

		log.Printf("%#v", command)
	}
}
