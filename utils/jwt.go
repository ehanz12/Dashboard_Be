package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//func untuk membuat JWT
func GenerateJWT(userID uint, email string) (string, error) {
	//cek di env ada JWT_SECRET
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set in environment variables")
	}
	// Implementation for generating JWT token
	//membuat claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp" : time.Now().Add(time.Hour * 24).Unix(), // token berlaku 24 jam
	}
	//membuat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//sign token
	return token.SignedString([]byte(secret))
}

//func untuk memvalidasi JWT
func VerifyJWT(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set in environment variables")
	}
	//memvalidasi token
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error){
		if _ , ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method") // jika metode signing tidak sesuai
		}
		return []byte(secret), nil
	})
}