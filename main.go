package main

import (
	"fmt"
	"net/http"

	"first-api-gin/config"
	"first-api-gin/responses"
	"first-api-gin/routes"

	"github.com/gin-gonic/gin"
)

func status(c *gin.Context) {
	c.IndentedJSON(
		http.StatusOK,
		responses.Response{
			Status: http.StatusOK,
			Msg:    "Hello",
			Data:   map[string]interface{}{"data": "World"},
		})
}

func main() {
	fmt.Println("Hello")
	router := gin.Default()
	config.ConnectDB()

	router.GET("/status", status)

	routes.CharacterRoute(router)

	router.Run("localhost:5000")
}
