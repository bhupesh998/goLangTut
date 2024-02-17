package myJWT

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/bhupesh998/go-csrf/db"
	"github.com/bhupesh998/go-csrf/db/models"
	
)

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath  = "keys/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func InitJWT() error {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

func CreateNewTokens(uuid string, role string) (authTokenString, refreshTokenString, csrfSecret string, err error) {

	csrfSecret, err = models.GenerateCSRFSecret()
	if err != nil {
		return
	}

	refreshTokenString, err = createRefreshTokenString(uuid, role, csrfSecret)
	if err != nil {
		return
	}

	authTokenString, err = createAuthTokenString(uuid, role, csrfSecret)
	if err != nil {
		return
	}

	return
}

func CheckAndRefreshTokens(oldAuthTokenString string, oldRefreshTokenString string, OldcsrfSecret string) (newAuthTokenString, newRefreshTokenString, newCsrfSecret string, err error) {
	if OldcsrfSecret == "" {
		log.Println("NO CSRF SECRET")
		err = errors.New("Unauthorized")
		return
	}

	authToken, err := jwt.ParseWithClaims(oldAuthTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	authTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
	if !ok {
		return
	}

	if OldcsrfSecret != authTokenClaims.CSRF {
		log.Printf("CSRF Token doesn't match jwt")
		err = errors.New("Unauthorized")
		return
	}

	if authToken.Valid {
		log.Println("Auth Token is Valid")

		newCsrfSecret = authTokenClaims.CSRF

		newRefreshTokenString, err = updateRefreshTokenExpiry(oldRefreshTokenString)

		newAuthTokenString = oldAuthTokenString
		return
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		log.Println("Auth token is not valid")
		if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			log.Println("Auth token is expired")

			newAuthTokenString, newCsrfSecret, err = updateAuthTokenString(oldRefreshTokenString, oldAuthTokenString)
			if err != nil {
				return
			}

			newRefreshTokenString, err = updateRefreshTokenExpiry(oldRefreshTokenString)
			if err != nil {
				return
			}

			newRefreshTokenString, err = updateRefreshTokenCsrf(newRefreshTokenString, newCsrfSecret)
			return
		} else {
			log.Println("Err in auth token")
			err = errors.New("Err in auth token")
			return
		}
	} else {
		log.Println("Err in auth token")
		err = errors.New("Err in auth token")
		return
	}

	err = errors.New("Unauthorized")
	return
}

func createAuthTokenString(uuid, role, csrfString string) (authTokenString string, err error) {
	authTokenExp := time.Now().Add(models.AuthTokenValidTime).Unix()
	authClaims := models.TokenClaims{
		jwt.StandardClaims{
			Subject:   uuid,
			ExpiresAt: authTokenExp,
		},
		role,
		csrfString,
	}

	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), authClaims)
	authTokenString, err = authJwt.SignedString(signKey)

	return
}

func createRefreshTokenString(uuid, role, csrfString string) (refreshTokenString string, err error) {
	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()
	refreshJti, err := db.StoreRefreshToken()
	if err != nil {
		return
	}

	refreshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        refreshJti,
			Subject:   uuid,
			ExpiresAt: refreshTokenExp,
		},
		role,
		csrfString,
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)
	refreshTokenString, err = refreshJwt.SignedString(signKey)

	return
}

func updateRefreshTokenExpiry(oldRefreshTokenString string) (newRefreshTokenString string, err error) {
	refreshToken, _ := jwt.ParseWithClaims(oldRefreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		return
	}

	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()

	refeshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        oldRefreshTokenClaims.StandardClaims.Id,
			Subject:   oldRefreshTokenClaims.StandardClaims.Subject,
			ExpiresAt: refreshTokenExp,
		},
		oldRefreshTokenClaims.Role,
		oldRefreshTokenClaims.CSRF,
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refeshClaims)

	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func updateAuthTokenString(refreshTokenString string, oldAuthTokenString string) (newAuthTokenString, csrfSecret string, err error) {

	refreshToken, _ := jwt.ParseWithClaims(refreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	refreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		err = errors.New("error reading the jwt claims")
		return
	}

	if db.CheckRefreshToken(refreshTokenClaims.StandardClaims.Id) {
		if refreshToken.Valid {
			authToken, _ := jwt.ParseWithClaims(oldAuthTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})

			authTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
			if !ok {
				err = errors.New("Error Fetching claims")
				return
			}

			csrfSecret, err = models.GenerateCSRFSecret()
			if err != nil {
				return
			}
			createAuthTokenString(authTokenClaims.StandardClaims.Subject, authTokenClaims.Role, csrfSecret)
			return
		} else {
			log.Println("Refresh Token has expired")
			db.DeleteRefreshToken(refreshTokenClaims.StandardClaims.Id)

			err = errors.New("Unauthorized")
			return
		}
	} else {
		log.Println("Refresh Token has Benn revoked")

		err = errors.New("Unauthorized")
		return
	}

}

func RevokeRefreshToken(refreshTokenString string) error {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return errors.New("Could Not Parse Refresh Token With Claims")
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		err = errors.New("Could Not Parse Refresh Token With Claims")

	}

	db.DeleteRefreshToken(refreshTokenClaims.StandardClaims.Id)
	return nil
}

func updateRefreshTokenCsrf(oldRefreshTokenString string, newCsrfString string) (newRefreshTokenString string, err error) {
	refreshToken, err := jwt.ParseWithClaims(oldRefreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		err = errors.New("Could Not Parse Refresh Token With Claims")

	}

	refeshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        oldRefreshTokenClaims.StandardClaims.Id,
			Subject:   oldRefreshTokenClaims.StandardClaims.Subject,
			ExpiresAt: oldRefreshTokenClaims.StandardClaims.ExpiresAt,
		},
		oldRefreshTokenClaims.Role,
		newCsrfString,
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refeshClaims)

	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return

}

func GrabUUID(authTokenString string) (string, error) {
	authToken, _ := jwt.ParseWithClaims(authTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return "", errors.New("Error Fetching claims")
	})

	authTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
	if !ok {
		return "", errors.New("Error Fetching claims")
	}

	return authTokenClaims.StandardClaims.Subject, nil
}
