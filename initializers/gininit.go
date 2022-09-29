package initializers

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/handlers"
)

var (
	r = gin.Default()
)

func CreateGinConnection() {
	handlers.CreateGinHandlers(r)

	ginErr := r.Run(":8080")

	if ginErr != nil {
		log.Printf("Error connecting to gin services %v", ginErr)
	} else {
		log.Println("Connected to gin services")
	}
}
