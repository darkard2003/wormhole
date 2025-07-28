package jwtservice

import (
	"fmt"
	"time"

	"github.com/darkard2003/wormhole/internals/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	AccessSecrectKey  string
	RefreshSecrectKey string
}

func NewJWTService(accessSecrectKey, refreshSecrectKey string) *JWTService {
	service := &JWTService{
		accessSecrectKey,
		refreshSecrectKey,
	}
	return service
}

func (s *JWTService) GenerateToken(userID int, userName string) (string, string, error) {
	exp := time.Now().Add(time.Hour)
	claims := &models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wormhole",
			ID:        uuid.New().String(),
		},
		UserID:    userID,
		UserName:  userName,
		TokenType: models.TOKEN_TYPE_ACCESS,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.AccessSecrectKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, exp.Format(time.RFC3339), nil
}

func (s *JWTService) GenerateRefereshToken(userID int, userName string) (string, string, error) {
	exp := time.Now().Add(time.Hour * 24 * 30)

	claims := &models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wormhole",
			ID:        uuid.New().String(),
		},
		UserID:    userID,
		UserName:  userName,
		TokenType: models.TOKEN_TYPE_REFRESH,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.RefreshSecrectKey))
	if err != nil {
		return "", "", err
	}
	return tokenString, exp.Format(time.RFC3339), nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.AccessSecrectKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.RefreshSecrectKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
