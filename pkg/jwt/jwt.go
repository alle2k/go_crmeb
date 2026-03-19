package jwt

import (
	"boilerplate/config"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte(config.AppConfig.Token.Secret)

type jwtClaims struct {
	UserID uint
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	claims := jwtClaims{
		UserID: userID,
	}
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Minute * time.Duration(config.AppConfig.Token.ExpireTime)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		log.Printf("Unable to parse token: %#v, err: %#v \n", tokenString, err)
		return 0, err
	}
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		log.Printf("Invalid token: %#v \n", token)
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}
