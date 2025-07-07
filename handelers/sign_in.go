package handelers

import (
	"crypto"
	"encoding/hex"

	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/darkard2003/wormhole/services/jwtservice"
	"github.com/gin-gonic/gin"
)

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignInHandeler(ctx *gin.Context) {
	var input SignInInput
	ctx.ShouldBindJSON(&input)

	jwt := jwtservice.GetJWTService()
	db := dbservice.GetDBService()

	user, err := db.GetUserByUsername(input.Username)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Database error"})
		return
	}

	if user == nil {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}

	hash := crypto.SHA256.New()
	hash.Write([]byte(input.Password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	if user.Password != hashedPassword {
		ctx.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	jwtToken, err := jwt.GenerateToken(input.Username)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(200, gin.H{
		"token":   jwtToken,
		"message": "Sign in successful",
	})
}
