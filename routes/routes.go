package routes

import (
	"api/services"
	temp "api/templates"
	"encoding/json"
	"net/http"
)

func Routes() {

	fs := http.FileServer(http.Dir("./templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/mainMenu", http.StatusSeeOther)
	})

	http.HandleFunc("/mainMenu", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "mainMenu", nil)
	})

	http.HandleFunc("/collection", func(w http.ResponseWriter, r *http.Request) {
		// Appel au service Marvel
		data, err := services.FetchMarvelCharacters()
		if err != nil {
			http.Error(w, "Erreur lors de l'appel à l'API Marvel", http.StatusInternalServerError)
			return
		}

		// Décoder les données JSON si besoin
		var characters map[string]interface{}
		err = json.Unmarshal([]byte(data), &characters)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture des données Marvel", http.StatusInternalServerError)
			return
		}

		// Envoyer les données à la page HTML
		temp.Render(w, "collection", characters)
	})

	http.HandleFunc("/ressource", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "ressource", nil)
	})

	http.HandleFunc("/categorie", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "categorie", nil)
	})

	http.HandleFunc("/favoris", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "favoris", nil)
	})

	http.HandleFunc("/recherche", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "recherche", nil)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "about", nil)
	})
}
