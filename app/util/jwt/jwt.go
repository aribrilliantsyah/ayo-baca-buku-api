package jwt

import (
	"ayo-baca-buku/app/util/logger"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	UID      int64  `json:"uid"`
	Username string `json:"username"`
}

func GenerateToken(UID string, userName string) (string, error) {
	logger := logger.GetLogger()

	jwtKey := []byte(viper.GetString("JWT_SECRET"))
	if len(jwtKey) == 0 {
		logger.Fatal("JWT_SECRET is not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      UID,
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (int64, error) {
	logger := logger.GetLogger()

	jwtKey := []byte(viper.GetString("JWT_SECRET"))
	if len(jwtKey) == 0 {
		logger.Fatal("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return 0, errors.New("token expired")
		}
		if userID, ok := claims["userid"].(float64); ok {
			return int64(userID), nil
		}
		return 0, errors.New("userid claim is missing or not a float64")
	}

	return 0, errors.New("invalid token")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DecodeToken(tokenString string) (*TokenClaims, error) {
	logger := logger.GetLogger()

	jwtKey := []byte(viper.GetString("JWT_SECRET"))
	if len(jwtKey) == 0 {
		logger.Fatal("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid := int64(claims["uid"].(float64))
		username := claims["username"].(string)
		return &TokenClaims{
			UID:      uid,
			Username: username,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func GetUserInfo(c *fiber.Ctx) (*TokenClaims, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}
	tokenString := authHeader[len("Bearer "):]
	tokenClaims, err := DecodeToken(tokenString)
	if err != nil {
		return nil, err
	}

	return tokenClaims, nil
}
