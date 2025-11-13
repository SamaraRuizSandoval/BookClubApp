package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/SamaraRuizSandoval/BookClubApp/docs"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/api"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/middleware"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/SamaraRuizSandoval/BookClubApp/migrations"
	"github.com/gin-gonic/gin"
)

// Resourses that we can use through our application
type Application struct {
	Logger         *log.Logger
	DB             *sql.DB
	Middleware     middleware.UserMiddleware
	BookHandler    *api.BookHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	CommentHandler *api.ChapterCommentHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	bookStore := store.NewPostgresBookStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	commentStore := store.NewPostgresChapterCommentStore(pgDB)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	bookHandler := api.NewBookHandler(bookStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	commentHandler := api.NewChapterCommentHandler(commentStore, logger)

	app := &Application{
		Logger:         logger,
		DB:             pgDB,
		Middleware:     middlewareHandler,
		BookHandler:    bookHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		CommentHandler: commentHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
