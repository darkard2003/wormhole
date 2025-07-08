package jwtservice

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	SecretKey string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func NewJWTService() *JWTService {
	service := &JWTService{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
	return service
}

func (s *JWTService) GenerateToken(userID int, userName string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 30)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wormhole",
			ID:        uuid.New().String(),
		},
		UserID:   userID,
		UserName: userName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
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
	return claims, nil
}
