package token

import (
	"context"
	"log"
	"os"
	"time"
	"github.com/bhupesh998/ecommproject/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

var userData *mongo.Collection = database.UserData(database.Client, "Users")

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

func TokenGenerator(email, firstname, lastname, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refeshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil{
		return "","", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refeshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil{
		log.Panic(err)
		return 
	}

	return token, refreshToken , err
}

func ValidateToken(signedToken string)(claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token)(interface{}, error){
		return []byte(SECRET_KEY), nil
	})

	if err!=nil{
		msg = err.Error()
		return 
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok{
		msg = "The Token is Invalid "
		return 
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = "The Token is Already expired"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)

	
	var updatedObj primitive.D

	updatedObj = append(updatedObj , bson.E{Key: "token", Value: signedToken})
	updatedObj = append(updatedObj , bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedObj = append(updatedObj , bson.E{Key: "updated_at", Value: updated_at})

	upsert := true

	filter := bson.M{"userid": userId}
	opt := options.UpdateOptions{
		Upsert : &upsert,
	}

	_, err := userData.UpdateOne(ctx,filter, bson.D{
		{Key: "$set", Value : updatedObj}, 
		}, &opt )

		defer cancel()

		if err !=nil {
			log.Panic(err)
			return
		}
}
