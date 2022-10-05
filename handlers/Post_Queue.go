package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/models"
)

func Queue(c *gin.Context) {
	var command models.Queue

	//bind it to the command struct
	err := c.BindJSON(&command)
	if err != nil {
		c.JSON(406, gin.H{
			"message": err,
		})
	}

	//check to make sure that all values are present
	if command.PlayerID != "" && command.Gamemode != "" && command.Primary != "" && command.Secondary != "" {
		c.JSON(200, gin.H{
			"message": "Data Received!",
		})
		log.Printf("%#v", command)
	} else {
		c.JSON(206, gin.H{
			"message": err,
		})
	}
}
