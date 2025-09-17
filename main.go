package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/app"
	"github.com/SamaraRuizSandoval/BookClubApp/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 5000, "go backend server port")
	flag.Parse() // Future note, use os.Getenv

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
