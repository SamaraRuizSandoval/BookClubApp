package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/SamaraRuizSandoval/BookClubApp/internal/api"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/routes"
)

// @title           BookClubApp
// @version         1.0
// @description     The BookClubApp to manage, share, and comment your favorite books. The goal is to create a space where you can interact and express your ideas and though as you go through the chapters of the books you are reading.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes        https
// @host      https://bookclub-backend.redwater-26f8bbd2.centralus.azurecontainerapps.io/
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	var port int
	flag.IntVar(&port, "port", 5000, "go backend server port")
	flag.Parse() // Future note, use os.Getenv

	// Azure Container Apps uses PORT env var
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := app.DB.Close(); err != nil {
			app.Logger.Printf("failed to close DB: %v", err)
		}
	}()

	r := routes.SetupRouter(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Printf("Running backend server on port %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("server error: %v", err)
	}
}
