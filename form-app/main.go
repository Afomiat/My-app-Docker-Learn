package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents the user data structure
type User struct {
	UserID    int    `bson:"userid" json:"userid"`
	Name      string `bson:"name" json:"name"`
	Email     string `bson:"email" json:"email"`
	Interests string `bson:"interests" json:"interests"`
}

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Serve static files (HTML, JS, CSS)
	r.Static("/static", "./static")

	// Serve HTML file at root
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Helper function to get MongoDB collection
	getCollection := func() (*mongo.Client, *mongo.Collection, context.Context, context.CancelFunc, error) {
		clientOptions := options.Client().ApplyURI("mongodb://admin:password@localhost:27017")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			cancel()
			return nil, nil, nil, nil, err
		}

		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			cancel()
			return nil, nil, nil, nil, err
		}

		collection := client.Database("user-account").Collection("users")
		return client, collection, ctx, cancel, nil
	}

	// MongoDB Connection endpoint matching the image
	r.GET("/get-profile", func(c *gin.Context) {
		client, collection, ctx, cancel, err := getCollection()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		defer client.Disconnect(ctx)

		// Find a document
		var result bson.M
		err = collection.FindOne(ctx, bson.M{"userid": 1}).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// New endpoint to update user profile
	r.POST("/update-profile", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Force userID to 1 for this single-user demo
		user.UserID = 1

		client, collection, ctx, cancel, err := getCollection()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		defer client.Disconnect(ctx)

		// Update or Insert (Upsert)
		opts := options.Update().SetUpsert(true)
		filter := bson.M{"userid": 1}
		update := bson.M{"$set": user}

		_, err = collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "form-app",
		})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
