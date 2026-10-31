package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/oscarracuna/portfolio/pkg/config"
	"github.com/oscarracuna/portfolio/pkg/handlers"
	"github.com/oscarracuna/portfolio/pkg/render"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Unable to load .env file.")
	}

	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		portNumber = ":8080"
	}

	if !strings.Contains(portNumber, ":") {
		portNumber = ":" + portNumber
	}

	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Unable to create template cache.")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("App has started on port %s\n", portNumber)

	srv := &http.Server{
		Addr:              portNumber,
		Handler:           routes(&app),
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
