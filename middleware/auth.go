package middleware

import (
	"net/http"

	"github.com/bhanupbalusu/gocomboums_v4/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		claims, err := service.VerifyAndExtractClaims(token, []byte("key"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// You can also store user information from the token in the context if needed
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
