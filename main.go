package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/AGX18/goWorkoutAPI/internal/app"
	"github.com/AGX18/goWorkoutAPI/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Logger.Printf("Starting application... Server started on :%d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("Error starting server: %v", err)
	}

}
