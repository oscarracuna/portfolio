package render

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oscarracuna/portfolio/pkg/config"
	"github.com/oscarracuna/portfolio/pkg/models"
)

func TestRenderEscapesHTML(t *testing.T) {
	previous := app
	t.Cleanup(func() { app = previous })
	page := template.Must(template.New("page").Parse(`<p>{{.StringMap.message}}</p>`))
	NewTemplates(&config.AppConfig{UseCache: true, TemplateCache: map[string]*template.Template{"page": page}})
	response := httptest.NewRecorder()
	RenderTemplate(response, "page", &models.TemplateData{StringMap: map[string]string{"message": "<script>alert(1)</script>"}})
	if response.Code != http.StatusOK || strings.Contains(response.Body.String(), "<script>") {
		t.Fatalf("unsafe or unsuccessful response: %d %s", response.Code, response.Body.String())
	}
}

func TestRenderFailureReturns500WithoutPartialPage(t *testing.T) {
	previous := app
	t.Cleanup(func() { app = previous })
	broken := template.Must(template.New("broken").Parse(`partial page {{.MissingField}}`))
	NewTemplates(&config.AppConfig{UseCache: true, TemplateCache: map[string]*template.Template{"broken": broken}})
	for _, name := range []string{"missing", "broken"} {
		response := httptest.NewRecorder()
		RenderTemplate(response, name, &models.TemplateData{})
		if response.Code != http.StatusInternalServerError {
			t.Errorf("%s: got %d", name, response.Code)
		}
		if strings.Contains(response.Body.String(), "partial page") {
			t.Error("sent a partial page on failure")
		}
	}
}
