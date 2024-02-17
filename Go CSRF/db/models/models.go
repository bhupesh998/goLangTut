package models

import (
	"time"
	"github.com/bhupesh998/go-csrf/randomstrings"
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct{
	Username, PasswordHash, Role string
}

type TokenClaims struct{
	jwt.StandardClaims
	Role string `json:"role"`
	CSRF string `json:"csrf"`
}

const RefreshTokenValidTime = time.Hour * 72
const AuthTokenValidTime = time.Minute * 15

func GenerateCSRFSecret()(string, error){
	return randomstrings.GenerateRandomString(32)
}

