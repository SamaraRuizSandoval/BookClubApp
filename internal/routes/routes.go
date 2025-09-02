package routes

import (
	"net/http"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *app.Application) *gin.Engine {
	r := gin.Default()

	r.GET("/health", app.HealthCheck)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})
	return r
}
