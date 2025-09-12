package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/gin-gonic/gin"
)

// Resourses that we can use through our application
type Application struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &Application{
		Logger: logger,
		DB:     pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
