package middleware

import (
	"net/http"
	"strings"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

func (um *UserMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set("user", store.AnonymusUser)
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		token := parts[1]
		user, err := um.UserStore.GetUserToken("authentication", token)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func (um *UserMiddleware) RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userValue, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in"})
			c.Abort()
			return
		}

		user, ok := userValue.(*store.User)
		if !ok || user.IsAnonymus() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in"})
			c.Abort()
			return
		}

		c.Next()
	}
}
