package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type Event struct {
	ID          string `json:"ID,omitempty" bson:"ID,omitempty"`
	Title       string `json:"Title,omitempty" bson:"Title,omitempty"`
	Description string `json:"Description,omitempty" bson:"Description,omitempty"`
}

func HomeLink(c *gin.Context) {
	c.String(http.StatusOK, "Welcome home!")
}

func CreateEvent(c *gin.Context) {
	var newEvent Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	insertResult, err := collection.InsertOne(context.TODO(), newEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insertResult)
}

func GetAllEvents(c *gin.Context) {
	var events []Event
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongodb:27017"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection = client.Database("synonyms").Collection("events")

	r := gin.Default()
	r.GET("/", HomeLink)
	r.POST("/event", CreateEvent)
	r.GET("/events", GetAllEvents)

	r.Run(":6666")
}
