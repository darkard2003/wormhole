package authhandelers

import (
	"log"

	"github.com/darkard2003/wormhole/interfaces"
	"github.com/darkard2003/wormhole/services/jwtservice"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignInHandlerHandler(db interfaces.DBInterface) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var input SignInInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		jwt := jwtservice.GetJWTService()

		user, err := db.GetUserByUsername(input.Username)
		if err != nil {
			log.Println("Error fetching user:", err)
			ctx.JSON(500, gin.H{"error": "Database error"})
			return
		}

		if user == nil {
			ctx.JSON(404, gin.H{"error": "Invalid username or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

		if err != nil {
			log.Println("Password mismatch:", err)
			ctx.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		jwtToken, err := jwt.GenerateToken(user.Id, user.Username)
		if err != nil {
			log.Println("Error generating JWT token:", err)
			ctx.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		ctx.JSON(200, gin.H{
			"token":   jwtToken,
			"message": "Sign in successful",
		})
	}
}
