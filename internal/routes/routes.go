package routes

import (
	"net/http"

	_ "github.com/SamaraRuizSandoval/BookClubApp/docs"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(app *app.Application) *gin.Engine {
	r := gin.Default()

	//"http://localhost:5173"
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://zealous-wave-0844b9e0f.2.azurestaticapps.net"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", app.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adminAuth := r.Group("/")
	adminAuth.Use(app.Middleware.AuthMiddleware(), app.Middleware.RequireAdmin())
	{
		adminAuth.POST("/books", app.BookHandler.HandleAddBook)
		adminAuth.POST("/admins", app.UserHandler.RegisterAdminAccount)
	}

	auth := r.Group("/")
	auth.Use(app.Middleware.AuthMiddleware(), app.Middleware.RequireUser())
	{
		auth.GET("/me", app.UserHandler.GetMe)
		auth.PUT("/books/:id", app.BookHandler.HandleUpdateBookByID)
		auth.DELETE("/books/:id", app.BookHandler.HandleDeleteBookByID)

		auth.POST("/chapters/:chapter_id/comments", app.CommentHandler.HandleAddComment)
		auth.PUT("/chapters/:chapter_id/comments/:id", app.CommentHandler.HandleUpdateComment)
		auth.DELETE("/chapters/:chapter_id/comments/:id", app.CommentHandler.HandleDeleteCommentById)
		auth.POST("/users/:user_id/books", app.UserBooksHandler.HandleAddUserBook)
		auth.GET("/users/:user_id/books", app.UserBooksHandler.HandleGetUserBooks)
		auth.PATCH("/user-books/:id", app.UserBooksHandler.HandleUpdateUserBook)
		auth.DELETE("/user-books/:id", app.UserBooksHandler.HandleDeleteUserBook)
	}

	r.GET("/books/:id", app.BookHandler.HandleGetBookByID)
	r.GET("/books", app.BookHandler.HandleGetAllBooks)

	r.GET("/chapters/:chapter_id/comments/", app.CommentHandler.HandleGetCommentsByChapterID)
	r.GET("/chapters/:chapter_id/comments/:id", app.CommentHandler.HandleGetCommentById)

	r.GET("/users", app.UserHandler.HandleGetUserByUsername)
	r.POST("/users", app.UserHandler.RegisterUser)
	r.POST("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})
	return r
}
