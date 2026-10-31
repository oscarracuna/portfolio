package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/oscarracuna/portfolio/pkg/config"
	"github.com/oscarracuna/portfolio/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		var err error
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Printf("load templates: %v", err)
			http.Error(w, "Unable to render page", http.StatusInternalServerError)
			return
		}
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Printf("template not found: %s", tmpl)
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	if err := t.Execute(buf, td); err != nil {
		log.Printf("execute template %s: %v", tmpl, err)
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		return
	}

	// render the template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
