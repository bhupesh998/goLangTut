package helpers

import (
	"context"
	"fmt"
	"go-resturant-management/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct{
	Email string
	First_name string
	Last_name string
	Uid string
	jwt.StandardClaims
}

func GenerateAllTokens(email string , first_name string, last_name string, user_id string)(signedToken string , signedRefreshToken string, err error){
	
	claims := &SignedDetails{
		Email: email,
		First_name: first_name,
		Last_name: last_name,
		Uid: user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims:=&SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token , err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token , refreshToken , err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

		filter := bson.M{"user_id": userId}
		upsert := true
		opt := options.UpdateOptions{
				Upsert: &upsert,
		}

		_ , updateErr := userCollection.UpdateOne(ctx, filter, bson.D{
			{"$set" ,updateObj},
		}, &opt)

		defer cancel()

		if updateErr != nil {
			log.Panic(updateErr)
			return
		}

		return
	
}


func ValidateToken(signedToken string )(claims *SignedDetails, msg string){

	token , err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token)(interface{}, error){
			return []byte(SECRET_KEY), nil
		},
	)

	claims , ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg

}