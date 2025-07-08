package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/darkard2003/wormhole/services/jwtservice"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		var token string
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with 'Bearer '"})
		} else {
			token = strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
				c.Abort()
				return
			}
		}

		jwt := jwtservice.GetJWTService()

		claims, err := jwt.ValidateToken(token)
		if err != nil {
			log.Println("Error validating token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims are nil"})
			c.Abort()
			return
		}

		userId := claims.UserID
		userName := claims.UserName

		if userId == 0 || userName == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Set("userName", claims.UserName)
		c.Next()
	}
}
