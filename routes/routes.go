package routes

import (
	"api/services"
	temp "api/templates"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func InitRoutes(cacheService *services.CacheService) {
	marvelService := services.NewMarvelService()

	fs := http.FileServer(http.Dir("./templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/mainMenu", http.StatusSeeOther)
	})

	// Pages statiques
	http.HandleFunc("/mainMenu", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "mainMenu", nil)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "about", nil)
	})

	// Collection avec pagination
	http.HandleFunc("/collection", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pageNum := 1
		if p, err := strconv.Atoi(page); err == nil {
			pageNum = p
		}
		limit := 20
		offset := (pageNum - 1) * limit

		data, err := cacheService.GetCharacters(offset, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response map[string]interface{}
		json.Unmarshal(data, &response)

		total := int(response["data"].(map[string]interface{})["total"].(float64))
		totalPages := (total + limit - 1) / limit

		temp.Render(w, "collection", map[string]interface{}{
			"data":        response["data"],
			"currentPage": pageNum,
			"totalPages":  totalPages,
			"limit":       limit,
		})
	})

	// Recherche
	http.HandleFunc("/recherche", func(w http.ResponseWriter, r *http.Request) {
		data, err := marvelService.FetchMarvelCharacters(0, 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response map[string]interface{}
		json.Unmarshal([]byte(data), &response)

		temp.Render(w, "recherche", map[string]interface{}{
			"data":  response["data"],
			"query": "",
		})
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			// Si la requête est vide, rediriger vers la page de recherche
			http.Redirect(w, r, "/recherche", http.StatusSeeOther)
			return
		}

		data, err := marvelService.SearchCharacters(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response map[string]interface{}
		json.Unmarshal([]byte(data), &response)

		temp.Render(w, "recherche", map[string]interface{}{
			"data":  response["data"],
			"query": query,
		})
	})

	// Catégories
	http.HandleFunc("/categorie", func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("cat")
		page := r.URL.Query().Get("page")
		pageNum := 1
		if p, err := strconv.Atoi(page); err == nil {
			pageNum = p
		}
		limit := 20
		offset := (pageNum - 1) * limit

		if category == "" {
			temp.Render(w, "categorie", nil)
			return
		}

		params := marvelService.GetCategoryParams(category)
		data, err := marvelService.FetchMarvelCharactersWithParams(offset, limit, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response map[string]interface{}
		json.Unmarshal([]byte(data), &response)

		total := int(response["data"].(map[string]interface{})["total"].(float64))
		totalPages := (total + limit - 1) / limit

		temp.Render(w, "categorie", map[string]interface{}{
			"data":        response["data"],
			"currentPage": pageNum,
			"totalPages":  totalPages,
			"category":    category,
		})
	})

	// Favoris
	http.HandleFunc("/favoris", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// Recuperation de la liste des favoris

			// Recuperation des données des perso

			// Affichage des données des perso
			temp.Render(w, "favoris", nil)
		case "POST":
			// Récuperation de l'id

			// Ajouter cette id au ficher JSON

			// Recuperation de la liste des favoris

			// Recuperation des données des perso

			// Affichage des données des perso
			w.WriteHeader(http.StatusOK)
		}
	})

	// Route API pour récupérer les détails d'un personnage
	http.HandleFunc("/api/character/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid character ID", http.StatusBadRequest)
			return
		}

		id := parts[3]
		characterID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid character ID", http.StatusBadRequest)
			return
		}

		data, err := marvelService.FetchCharacterDetails(characterID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	})

	// Route pour les détails des personnages
	http.HandleFunc("/character/details/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/character/details/")

		// Convertir d'abord en float64 pour gérer la notation scientifique
		floatID, err := strconv.ParseFloat(path, 64)
		if err != nil {
			http.Error(w, "Invalid character ID", http.StatusBadRequest)
			return
		}

		// Convertir en int
		id := int(floatID)

		// Récupérer les détails du personnage
		characterData, err := marvelService.FetchCharacterDetails(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Récupérer les comics du personnage
		comicsData, err := marvelService.FetchCharacterComics(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Récupérer les séries du personnage
		seriesData, err := marvelService.FetchCharacterSeries(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var characterResponse, comicsResponse, seriesResponse map[string]interface{}
		json.Unmarshal([]byte(characterData), &characterResponse)
		json.Unmarshal([]byte(comicsData), &comicsResponse)
		json.Unmarshal([]byte(seriesData), &seriesResponse)

		character := characterResponse["data"].(map[string]interface{})["results"].([]interface{})[0]
		comics := comicsResponse["data"].(map[string]interface{})["results"].([]interface{})
		series := seriesResponse["data"].(map[string]interface{})["results"].([]interface{})

		temp.Render(w, "character_details", map[string]interface{}{
			"character": character,
			"comics":    comics,
			"series":    series,
		})
	})

	// Routes pour les ressources spécifiques
	http.HandleFunc("/api/comics/", func(w http.ResponseWriter, r *http.Request) {
		comicType := strings.TrimPrefix(r.URL.Path, "/api/comics/")
		offset := 0
		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			offset, _ = strconv.Atoi(offsetStr)
		}

		data, err := marvelService.FetchComicsByType(comicType, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response map[string]interface{}
		json.Unmarshal([]byte(data), &response)

		title := "Comics Marvel"
		switch comicType {
		case "popular":
			title = "Séries Populaires"
		case "avengers":
			title = "Comics Avengers"
		case "xmen":
			title = "Comics X-Men"
		}

		temp.Render(w, "resource_results", map[string]interface{}{
			"data":  response["data"],
			"title": title,
		})
	})

	http.HandleFunc("/api/team/", func(w http.ResponseWriter, r *http.Request) {
		team := strings.TrimPrefix(r.URL.Path, "/api/team/")
		offset := 0
		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			offset, _ = strconv.Atoi(offsetStr)
		}

		data, err := marvelService.FetchCharactersByTeam(team, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response map[string]interface{}
		json.Unmarshal([]byte(data), &response)

		title := "Équipe Marvel"
		switch team {
		case "avengers":
			title = "Les Avengers"
		case "xmen":
			title = "Les X-Men"
		}

		temp.Render(w, "resource_results", map[string]interface{}{
			"data":  response["data"],
			"title": title,
		})
	})

	http.HandleFunc("/api/universe/", func(w http.ResponseWriter, r *http.Request) {
		universe := strings.TrimPrefix(r.URL.Path, "/api/universe/")
		offset := 0
		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			offset, _ = strconv.Atoi(offsetStr)
		}

		data, err := marvelService.FetchUniverseCharacters(universe, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response map[string]interface{}
		json.Unmarshal([]byte(data), &response)

		title := "Univers Marvel"
		switch universe {
		case "616":
			title = "Terre-616"
		case "ultimate":
			title = "Ultimate Universe"
		}

		temp.Render(w, "resource_results", map[string]interface{}{
			"data":  response["data"],
			"title": title,
		})
	})

	// Après les autres routes statiques
	http.HandleFunc("/ressource", func(w http.ResponseWriter, r *http.Request) {
		temp.Render(w, "ressource", nil)
	})
}
