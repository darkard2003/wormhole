package authhandelers

import (
	"github.com/darkard2003/wormhole/internals/models"
	"github.com/darkard2003/wormhole/internals/server/middleware"
	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/gin-gonic/gin"
)

type UserDetailsResponse struct {
	Id       int     `json:"id"`
	Username string  `json:"username"`
	Email    *string `json:"email,omitempty"`
}

func UserDetails(db db.DBInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName, ok := middleware.GetUserName(c)
		if !ok {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := db.GetUserByUsername(userName)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(200, gin.H{"user": getUserResponse(user)})

	}
}

func getUserResponse(user *models.User) *UserDetailsResponse {
	return &UserDetailsResponse{
		user.Id,
		user.Username,
		user.Email,
	}
}
