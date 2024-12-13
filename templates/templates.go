package templates

import (
	"net/http"
	"text/template"
)

var templates *template.Template

// InitTemplates charge tous les templates du dossier "templates/"
func InitTemplates() {
	templates = template.Must(template.ParseGlob("templates/static/temp/*.html"))
}

// Render rend un template spécifique avec des données
func Render(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
