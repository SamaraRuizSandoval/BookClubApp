package middleware

import (
	"net/http"
	"strings"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

func GetUser(c *gin.Context) *store.User {
	userValue, ok := c.Get("user")
	if !ok {
		return nil
	}
	user, ok := userValue.(*store.User)
	if !ok || user.IsAnonymus() {
		return nil
	}
	return user
}

func (um *UserMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set("user", store.AnonymusUser)
			c.Next()
			return
		}

		parts := strings.Fields(strings.TrimSpace(authHeader))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
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
		if user.Role == RoleAdmin {
			c.Set("admin", user)
		}

		c.Next()
	}
}

func (um *UserMiddleware) RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (um *UserMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminValue, ok := c.Get("admin")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
			c.Abort()
			return
		}

		admin, ok := adminValue.(*store.User)
		if !ok || admin.IsAnonymus() {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
