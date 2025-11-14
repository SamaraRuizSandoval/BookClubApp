package routes

import (
	"net/http"

	_ "github.com/SamaraRuizSandoval/BookClubApp/docs"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(app *app.Application) *gin.Engine {
	r := gin.Default()

	r.GET("/health", app.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/")
	auth.Use(app.Middleware.AuthMiddleware(), app.Middleware.RequireUser())
	{
		auth.POST("/books", app.BookHandler.HandleAddBook)
		auth.PUT("/books/:id", app.BookHandler.HandleUpdateBookByID)
		auth.DELETE("/books/:id", app.BookHandler.HandleDeleteBookByID)

		auth.POST("/chapters/:chapter_id/comments", app.CommentHandler.HandleAddComment)
		auth.PUT("/chapters/:chapter_id/comments/:id", app.CommentHandler.HandleUpdateComment)
		auth.DELETE("/chapters/:chapter_id/comments/:id", app.CommentHandler.HandleDeleteCommentById)
	}

	r.GET("/books/:id", app.BookHandler.HandleGetBookByID)

	r.GET("/chapters/:chapter_id/comments/", app.CommentHandler.HandleGetCommentsByChapterID)
	r.GET("/chapters/:chapter_id/comments/:id", app.CommentHandler.HandleGetCommentById)

	r.GET("/users/:username", app.UserHandler.HandleGetUserByUsername)
	r.POST("/users", app.UserHandler.RegisterUser)
	r.POST("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})
	return r
}
