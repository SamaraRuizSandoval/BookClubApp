package routes

import (
	"net/http"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *app.Application) *gin.Engine {
	r := gin.Default()

	r.GET("/health", app.HealthCheck)
	r.GET("/books/:id", app.BookHandler.HandleGetBookByID)
	r.POST("/books", app.BookHandler.HandleAddBook)
	r.PUT("/books/:id", app.BookHandler.HandleUpdateBookByID)
	r.DELETE("/books/:id", app.BookHandler.HandleDeleteBookByID)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})
	return r
}
