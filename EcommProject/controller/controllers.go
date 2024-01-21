package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	generate "github.com/bhupesh998/ecommproject/token"
	"github.com/bhupesh998/ecommproject/database"
	"github.com/bhupesh998/ecommproject/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var prodCollection  *mongo.Collection = database.UserData(database.Client , "Products")
var userCollection  *mongo.Collection = database.ProductData(database.Client, "Users")
var Validate = validator.New()

func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err !=nil {
		log.Panic(err)
	}
	return string(bytes)
}

func verifyPassword(userPassword string, givenPassword string) (bool,string) {
	err:= bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Id or Password is Incorrect"
		valid = false
	}

	return valid , msg
}

func SignUp() gin.HandlerFunc{

	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err !=nil {
			c.JSON(http.StatusBadRequest , gin.H{"error": err.Error} )
			return
		}

		ValidationErr := Validate.Struct(user)
		if  ValidationErr !=nil{
			c.JSON(http.StatusBadRequest , gin.H{"error": ValidationErr} )
			return
		}

		count , err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err !=nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err})
			return
		}

		if count > 0{
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}

		count , err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err !=nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err})
			return
		}

		if count > 0{
			c.JSON(http.StatusBadRequest, gin.H{"error": "user with phone number already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At , _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		token , refreshToken,_ := generate.TokenGenerator(*user.Email , *user.First_Name, *user.Last_Name, user.UserId)

		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.Address_Details = make([]models.Address,0)
		user.UserCart = make([]models.ProductUser,0 )
		user.Order_Status = make([]models.Order, 0)

		_,err = userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user cannot be created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "Successfully Signed In")

	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
		var foundUser models.User
		var user models.User
		//whatever you receive from json will get binded to user variable
		if err :=c.BindJSON(&user); err !=nil{
			c.JSON(http.StatusBadRequest , gin.H{"error": err.Error} )
			return
			
		} 

		err := userCollection.FindOne(ctx , bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email or Password"})
			return
		}

		PasswordIsValid, msg := verifyPassword(*user.Password , *foundUser.Password)
		defer cancel()

		if !PasswordIsValid{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token , refreshToken, _ := generate.TokenGenerator(*foundUser.Email , *foundUser.First_Name, *foundUser.Last_Name,foundUser.UserId )
		defer cancel()
	
		generate.UpdateAllTokens(token , refreshToken,foundUser.UserId )

		c.JSON(http.StatusFound , foundUser)

	}
}


func ProductViewerAdmin() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var products models.Product
		if err := c.BindJSON(&products); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_Id = primitive.NewObjectID()
		_, anyErr := prodCollection.InsertOne(ctx,products)
		if anyErr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Inserted" })
			return 
		}

		defer cancel()
		c.JSON(http.StatusOK, "Successfully added ")
	}
}


func SearchProduct() gin.HandlerFunc{
	return func(c *gin.Context){
		var productList []models.Product

		var ctx,cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		cursor, err :=prodCollection.Find(ctx , bson.D{{}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went Wrong")
		}

		err = cursor.All(ctx , &productList)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		defer cursor.Close(ctx)

		if err := cursor.Err() ; err !=nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productList )
	}
}

func SearchProductByQuery() gin.HandlerFunc{
	return func(c *gin.Context){
		
		var searchProducts []models.Product
		queryParams := c.Query("name")

		if queryParams == ""{
			log.Println("query is empty")
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return 
		}

		var ctx,cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		searchQueryDb, err :=prodCollection.Find(ctx , bson.M{"product_name": bson.M{"$regex": queryParams}})

		if err !=nil {
			c.IndentedJSON(404, "Something went Wrong while fetching data")
			return
		}

		err = searchQueryDb.All(ctx , &searchProducts)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		defer searchQueryDb.Close(ctx)

		if err := searchQueryDb.Err() ; err !=nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchProducts )
	}
}