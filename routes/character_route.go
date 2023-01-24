package routes

import (
	"first-api-gin/controllers"

	"github.com/gin-gonic/gin"
)

func CharacterRoute(router *gin.Engine) {
	router.GET("/characters", controllers.GetCharacters())
	router.POST("/characters", controllers.CreateCharacter())
}
