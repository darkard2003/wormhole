package models

import "github.com/golang-jwt/jwt/v5"

type JWTTokenType string

const (
	TOKEN_TYPE_ACCESS  = "access"
	TOKEN_TYPE_REFRESH = "refresh"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID    int          `json:"user_id"`
	UserName  string       `json:"user_name"`
	TokenType JWTTokenType `json:"token_type"`
}
