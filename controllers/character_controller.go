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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			c.JSON(
				http.StatusBadRequest,
				responses.Response{
					Status: http.StatusBadRequest,
					Msg:    "error",
					Data:   map[string]interface{}{"data": err.Error()},
				})
			return
		}

		//validation with library
		if validationErr := validate.Struct(&character); validationErr != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.Response{
					Status: http.StatusBadRequest,
					Msg:    "error",
					Data:   map[string]interface{}{"data": validationErr.Error()},
				})
			return
		}

		newCharacter := models.Character{
			ID:        primitive.NewObjectID(),
			Name:      character.Name,
			Birthday:  character.Birthday,
			Dead:      character.Dead,
			Relevance: character.Relevance,
			Seasons:   character.Seasons,
		}

		result, err := characterCollection.InsertOne(ctx, newCharacter)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.Response{
					Status: http.StatusInternalServerError,
					Msg:    "error",
					Data:   map[string]interface{}{"data": err.Error()},
				})
			return
		}
		c.JSON(
			http.StatusCreated,
			responses.Response{
				Status: http.StatusCreated,
				Msg:    "success",
				Data:   map[string]interface{}{"data": result},
			})
	}
}

func GetCharacters() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var characters []models.Character

		results, err := characterCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.Response{
					Status: http.StatusInternalServerError,
					Msg:    "error",
					Data:   map[string]interface{}{"data": err.Error()},
				})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleCharacter models.Character
			if err = results.Decode(&singleCharacter); err != nil {
				c.JSON(
					http.StatusInternalServerError,
					responses.Response{
						Status: http.StatusInternalServerError,
						Msg:    "error",
						Data:   map[string]interface{}{"data": err.Error()},
					})
				return
			}
			characters = append(characters, singleCharacter)
		}

		c.JSON(
			http.StatusOK,
			responses.Response{
				Status: http.StatusOK,
				Msg:    "success",
				Data:   map[string]interface{}{"data": characters},
			})
	}
}

func GetOneCharacter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		characterId := c.Param("characterId")
		var character models.Character
		objId, _ := primitive.ObjectIDFromHex(characterId)

		err := characterCollection.FindOne(
			ctx,
			bson.M{"_id": objId},
		).Decode(&character)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.Response{
					Status: http.StatusInternalServerError,
					Msg:    "error",
					Data:   map[string]interface{}{"data": err.Error()},
				})
			return
		}

		c.JSON(
			http.StatusOK,
			responses.Response{
				Status: http.StatusOK,
				Msg:    "success",
				Data:   map[string]interface{}{"data": character},
			})
	}
}

func EditCharacter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		characterId := c.Param("characterId")
		var character models.Character
		objId, _ := primitive.ObjectIDFromHex(characterId)

		//validation body
		if err := c.BindJSON(&character); err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.Response{
					Status: http.StatusBadRequest,
					Msg:    "error",
					Data:   map[string]interface{}{"data": err.Error()},
				})
			return
		}

		updater := bson.M{"$set": bson.M{}}
		if character.Name != "" {
			updater["$set"].(bson.M)["name"] = character.Name
		}
		if character.Birthday != "" {
			updater["$set"].(bson.M)["birthday"] = character.Birthday
		}
		if character.Relevance != "" {
			updater["$set"].(bson.M)["relevance"] = character.Relevance
		}
		if character.Seasons != 0 {
			updater["$set"].(bson.M)["seasons"] = character.Seasons
		}

		result, err := characterCollection.UpdateOne(
			ctx,
			bson.M{"_id": objId},
			updater,
		)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.Response{
					Status: http.StatusBadRequest,
					Msg:    "error",
					Data:   map[string]interface{}{"data": err.Error()},
				})
			return
		}

		var updatedCharacter models.Character
		if result.MatchedCount == 1 {
			err := characterCollection.FindOne(
				ctx,
				bson.M{"_id": objId},
			).Decode(&updatedCharacter)
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					responses.Response{
						Status: http.StatusInternalServerError,
						Msg:    "error",
						Data:   map[string]interface{}{"data": err.Error()},
					})
				return
			}
		}

		c.JSON(
			http.StatusOK,
			responses.Response{
				Status: http.StatusOK,
				Msg:    "success",
				Data:   map[string]interface{}{"data": updatedCharacter},
			})
	}
}

func DeleteCharacter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		characterId := c.Param("characterId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(characterId)

		result, err := characterCollection.DeleteOne(
			ctx,
			bson.M{"_id": objId},
		)

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.Response{
					Status: http.StatusInternalServerError,
					Msg:    "error",
					Data:   map[string]interface{}{"data": err.Error()},
				})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(
				http.StatusNotFound,
				responses.Response{
					Status: http.StatusNotFound,
					Msg:    "error",
					Data:   map[string]interface{}{"data": "User with specified ID not found!"}},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			responses.Response{
				Status: http.StatusOK,
				Msg:    "deleted",
				Data:   map[string]interface{}{"data": result},
			})
	}
}
