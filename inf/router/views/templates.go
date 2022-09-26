package views

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var files embed.FS

var templates *template.Template

func init() {
	tmpl, err := template.New(``).ParseFS(files, `templates/*`)
	if err != nil {
		panic(err)
	}

	templates = tmpl
}

func Render(w http.ResponseWriter, name string, data any) {
	w.WriteHeader(http.StatusOK)

	err := templates.ExecuteTemplate(w, name+`.html`, data)
	if err != nil {
		panic(err)
	}
}
