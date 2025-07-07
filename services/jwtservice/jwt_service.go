package jwtservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	SecretKey string
}

type Claims struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
}

func NewJWTService() *JWTService {
	service := &JWTService{
		SecretKey: generateRandomSecretKey(),
	}
	return service
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not parse claims")
	}
	return &Claims{
		UserID: claim["user_id"].(string),
		Exp:    int64(claim["exp"].(float64)),
	}, nil
}

func generateRandomSecretKey() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()")
	secretKey := make([]rune, 32)
	for i := range secretKey {
		secretKey[i] = letters[rand.Intn(len(letters))]
	}
	return string(secretKey)
}
