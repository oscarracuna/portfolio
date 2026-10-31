package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/oscarracuna/portfolio/pkg/config"
	"github.com/oscarracuna/portfolio/pkg/handlers"
	"github.com/oscarracuna/portfolio/pkg/render"
)

// Exercise the real templates and asset routes to catch document-composition
// regressions such as rendering a second page just to include its styles.
func TestPortfolioRoutes(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	cache, err := render.CreateTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	app := &config.AppConfig{TemplateCache: cache, UseCache: true}
	handlers.NewHandlers(handlers.NewRepo(app))
	render.NewTemplates(app)
	router := routes(app)

	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/", nil))
	if page.Code != http.StatusOK {
		t.Fatalf("home returned %d", page.Code)
	}
	body := page.Body.String()
	for _, tag := range []string{"<!DOCTYPE html>", "<html ", "<body>", "</body>", "</html>"} {
		if strings.Count(body, tag) != 1 {
			t.Errorf("expected one %q", tag)
		}
	}
	if contentType := page.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("unexpected content type %q", contentType)
	}
	for _, asset := range []string{"/static/style.css", "/static/favicon.svg"} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, asset, nil))
		if response.Code != http.StatusOK {
			t.Errorf("%s returned %d", asset, response.Code)
		}
	}
}
