package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
	"github.com/bhupesh998/ecommproject/database"
	"github.com/bhupesh998/ecommproject/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("ProductID is Empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productID is Empty"))
			return
		}

		userId := c.Query("user_id")
		if userId == "" {
			log.Println("userId is Empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userId is Empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err )
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()
		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Added to Cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("ProductID is Empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productID is Empty"))
			return
		}

		userId := c.Query("user_id")
		if userId == "" {
			log.Println("userId is Empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userId is Empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()
		err = database.RemoveItemFromCart(ctx, app.prodCollection, app.userCollection, productID, userId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Removed Item From Cart")
	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id :=c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{ "error": "Invalid User Id"})
			c.Abort()
			return
		}

		userObjectId , err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500 , "Internal Server Error")
		}

		var ctx , cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userObjectId}}).Decode(&filledCart)


		if err !=nil{
			log.Println(err)
			c.IndentedJSON(500, "Not Found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value : bson.D{primitive.E{Key: "_id", Value : userObjectId}}}}
		grouping := bson.D{{Key: "$group", Value: primitive.E{Key:"_id",Value:"$id"}}} 
		unwind := bson.D{{Key: "$unwind", Value : bson.D{ primitive.E{Key: "path", Value: "$usercart"}, {Key: "Total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}} }}

		pointcursor , err := userCollection.Aggregate(ctx,mongo.Pipeline{filter_match, unwind, grouping})

		if err != nil{
			log.Println(err)
		}

		var listing []bson.M
		if err = pointcursor.All(ctx,&listing); err !=nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, json := range listing{
			c.IndentedJSON(200, json["Total"])
			c.IndentedJSON(200, filledCart.UserCart)
		}

		ctx.Done()





	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("user_id")
		if userId == "" {
			log.Panicln("userId is Empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userId is Empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Order Placed")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("ProductID is Empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("productID is Empty"))
			return
		}

		userId := c.Query("user_id")
		if userId == "" {
			log.Println("userId is Empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userId is Empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()
		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Order Placed Successfully")
	}
}
