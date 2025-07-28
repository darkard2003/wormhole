package authhandelers

import (
	"log"
	"net/http"

	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/darkard2003/wormhole/internals/services/jwtservice"
	"github.com/darkard2003/wormhole/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignInHandlerHandler(db db.DBInterface) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var input SignInInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		jwt := jwtservice.GetJWTService()

		user, err := db.GetUserByUsername(input.Username)
		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}

		if user == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

		if err != nil {
			log.Println("Password mismatch:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
			return
		}

		jwtToken, exp, err := jwt.GenerateToken(user.Id, user.Username)
		if err != nil {
			log.Println("Error generating JWT token:", err)
			ctx.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		refereshToken, refereshExp, err := jwt.GenerateRefereshToken(user.Id, user.Username)

		if err != nil {
			log.Println("Error generating JWT token:", err)
			ctx.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		ctx.JSON(200, gin.H{
			"token":              jwtToken,
			"expires_at":         exp,
			"refresh_token":      refereshToken,
			"refresh_expires_at": refereshExp,
			"message":            "Sign in successful",
		})
	}
}
