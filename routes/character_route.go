package routes

import (
	"first-api-gin/controllers"

	"github.com/gin-gonic/gin"
)

func CharacterRoute(router *gin.Engine) {
	router.GET("/characters", controllers.GetCharacters())
	router.GET("/characters/:characterId", controllers.GetOneCharacter())
	router.POST("/characters", controllers.CreateCharacter())
	router.PUT("/characters/:characterId", controllers.EditCharacter())
	router.DELETE("/characters/:characterId", controllers.DeleteCharacter())
}
