package routes

import (
	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *app.Application) *gin.Engine {
	r := gin.Default()

	r.GET("/health", app.HealthCheck)
	return r
}
