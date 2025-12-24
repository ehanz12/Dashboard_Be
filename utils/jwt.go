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
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		return "", errors.New("JWT_ACCESS_SECRET is not set in environment variables")
	}
	// Implementation for generating JWT token
	//membuat claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp" : time.Now().Add(15 * time.Minute).Unix(), // token berlaku 24 jam
	}
	//membuat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//sign token
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return "", errors.New("Error no found JWT_REFRESH_SECRET")
	}

	claims := jwt.MapClaims{
		"user_id" : userID,
		"exp" : time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

//func untuk memvalidasi JWT
func VerifyRefreshToken(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_REFRESH_SECRET is not set in environment variables")
	}
	//memvalidasi token
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error){
		if _ , ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method") // jika metode signing tidak sesuai
		}
		return []byte(secret), nil
	})
}


func VerifyAccessToken(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		return nil, errors.New("error bro not define")
	}
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error bor")
		}
		return []byte(secret), nil
	})
}