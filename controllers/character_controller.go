package controllers

import (
	"context"
	"first-api-gin/config"
	"first-api-gin/models"
	"first-api-gin/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var characterCollection *mongo.Collection = config.GetCollection(config.DB, "characters")
var validate = validator.New()

func CreateCharacter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var character models.Character
		defer cancel()

		//validacion body
		if err := c.BindJSON(&character); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Msg: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if validationErr := validate.Struct(&character); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Msg: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		}

		newCharacter := models.Character{
			// ID:        primitive.NewObjectID(),
			Name:      character.Name,
			Birthday:  character.Birthday,
			Dead:      character.Dead,
			Relevance: character.Relevance,
			Seasons:   character.Seasons,
		}

		result, err := characterCollection.InsertOne(ctx, newCharacter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Msg: "error", Data: map[string]interface{}{"data": err.Error()}})
		}
		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Msg: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetCharacters() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var characters []models.Character

		results, err := characterCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Msg: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleCharacter models.Character
			if err = results.Decode(&singleCharacter); err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Msg: "error", Data: map[string]interface{}{"data": err.Error()}})
			}
			characters = append(characters, singleCharacter)
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Msg: "success", Data: map[string]interface{}{"data": characters}})
	}
}
