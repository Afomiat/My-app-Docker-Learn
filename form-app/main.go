package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	UserID    int    `bson:"userid" json:"userid"`
	Name      string `bson:"name" json:"name"`
	Email     string `bson:"email" json:"email"`
	Interests string `bson:"interests" json:"interests"`
}

// getCollection is now a helper function that can be used by both main and tests
func getCollection() (*mongo.Client, *mongo.Collection, context.Context, context.CancelFunc, error) {
	mongoHost := os.Getenv("MONGO_URL")
	if mongoHost == "" {
		mongoHost = "localhost:27017"
	}

	uri := "mongodb://admin:password@" + mongoHost
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, nil, nil, nil, err
	}

	collection := client.Database("user-account").Collection("users")
	return client, collection, ctx, cancel, nil
}

// SetupRouter contains all your API logic
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.POST("/update-profile", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.UserID = 1

		client, collection, ctx, cancel, err := getCollection()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		defer client.Disconnect(ctx)

		opts := options.Update().SetUpsert(true)
		_, err = collection.UpdateOne(ctx, bson.M{"userid": 1}, bson.M{"$set": user}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	return r
}

func main() {
	r := SetupRouter()
	r.Run(":8080")
}