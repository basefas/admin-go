package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Claims struct {
	UID uint64
	jwt.RegisteredClaims
}

func GenerateToken(userID uint64) (string, error) {
	jwtSecret := []byte(viper.GetString("app.jwtSecret"))
	now := time.Now()
	expireTime := now.Add(time.Second * time.Duration(viper.GetInt("app.jwtTimeout")))
	claims := Claims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    viper.GetString("app.name"),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

func ParseToken(tokenString string) (*Claims, error) {
	jwtSecret := []byte(viper.GetString("app.jwtSecret"))

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
