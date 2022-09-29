package handlers

import "github.com/gin-gonic/gin"

func CreateGinHandlers(r *gin.Engine) {

	r.GET("/ping", Ping)
	r.POST("/queue", Queue)

}
