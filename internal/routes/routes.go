package routes

import (
	"net/http"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *app.Application) *gin.Engine {
	r := gin.Default()

	r.GET("/health", app.HealthCheck)

	auth := r.Group("/")
	auth.Use(app.Middleware.AuthMiddleware(), app.Middleware.RequireUser())
	{
		auth.POST("/books", app.BookHandler.HandleAddBook)
		auth.PUT("/books/:id", app.BookHandler.HandleUpdateBookByID)
		auth.DELETE("/books/:id", app.BookHandler.HandleDeleteBookByID)
	}

	r.GET("/books/:id", app.BookHandler.HandleGetBookByID)

	r.GET("/users/:username", app.UserHandler.HandleGetUserByUsername)
	r.POST("/users", app.UserHandler.HandleRegisterUser)
	r.POST("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})
	return r
}
