package handlers

import "github.com/gin-gonic/gin"

func CreateGinHandlers(r *gin.Engine) {

	r.POST("/ping", Ping)

}
