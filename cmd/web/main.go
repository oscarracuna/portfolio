package main

import (
	"fmt"
	"log"
	"net/http"
  "os"

	"github.com/joho/godotenv"

  "github.com/oscarracuna/portfolio/pkg/config"
	"github.com/oscarracuna/portfolio/pkg/handlers"
	"github.com/oscarracuna/portfolio/pkg/render"
)

// TODO: add terminal-looking flexbox for the whoami


func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Unable to load .env file.")
  }

  portNumber := os.Getenv("PORT")
  if portNumber == "" {
    log.Fatal("Unable to find a port inside .env file.")
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
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
