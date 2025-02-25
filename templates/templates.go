package templates

import (
	"net/http"
	"text/template"
)

var templates *template.Template

// Fonctions utilitaires pour les templates
var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"div": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return (a + b - 1) / b // Division avec arrondi supérieur
	},
	"le": func(a, b int) bool {
		return a <= b
	},
	"ge": func(a, b int) bool {
		return a >= b
	},
}

// InitTemplates charge tous les templates du dossier "templates/"
func InitTemplates() {
	templates = template.Must(template.New("").Funcs(templateFuncs).ParseGlob("templates/static/temp/*.html"))
}

// Render rend un template spécifique avec des données
func Render(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
