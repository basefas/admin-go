package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var jwtSecret []byte

type Claims struct {
	UID uint64
	jwt.StandardClaims
}

func GenerateToken(userID uint64) (string, error) {
	jwtSecret = []byte(viper.GetString("app.jwtSecret"))
	now := time.Now()
	expireTime := now.Add(time.Second * time.Duration(viper.GetInt("app.jwtTimeout")))
	claims := Claims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    viper.GetString("app.name"),
			IssuedAt:  now.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(tokenString string) (*Claims, error) {
	jwtSecret = []byte(viper.GetString("app.jwtSecret"))

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	claims := token.Claims.(*Claims)
	return claims, nil
}

func GetUID(tokenString string) (uint64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UID, nil
}
