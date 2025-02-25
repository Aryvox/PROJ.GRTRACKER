package controllers

import (
	"api/services"
	temp "api/templates"
	"net/http"
	"strconv"
)

type CharacterController struct {
	service *services.MarvelService
}

func NewCharacterController() *CharacterController {
	return &CharacterController{
		service: services.NewMarvelService(),
	}
}

func (c *CharacterController) HandleFavorite(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		characterID := r.FormValue("characterId")
		id, _ := strconv.Atoi(characterID)

		data, err := c.service.FetchCharacterDetails(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		characters, err := c.service.ProcessCharacters([]byte(data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(characters) == 0 {
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}

		character := characters[0]
		favorite := services.Favorite{
			ID:          character.ID,
			Name:        character.Name,
			Description: character.Description,
			Thumbnail: struct {
				Path      string `json:"path"`
				Extension string `json:"extension"`
			}{
				Path:      character.Thumbnail.Path,
				Extension: character.Thumbnail.Extension,
			},
		}

		if err := c.service.SaveFavorite(favorite); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)

	case "DELETE":
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if err := c.service.RemoveFavorite(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/favoris", http.StatusSeeOther)
	}
}

func (c *CharacterController) HandleCollection(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de pagination
	page := r.URL.Query().Get("page")
	pageNum := 1
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			pageNum = p
		}
	}
	limit := 20
	offset := (pageNum - 1) * limit

	data, err := c.service.FetchMarvelCharacters(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	characters, err := c.service.ProcessCharacters([]byte(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp.Render(w, "collection", map[string]interface{}{
		"data": map[string]interface{}{
			"results": characters,
			"total":   len(characters),
		},
		"currentPage": pageNum,
		"limit":       limit,
	})
}

func (c *CharacterController) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	data, err := c.service.SearchCharacters(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	characters, err := c.service.ProcessCharacters([]byte(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp.Render(w, "recherche", map[string]interface{}{
		"data":  map[string]interface{}{"results": characters},
		"query": query,
	})
}

func (c *CharacterController) HandleCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("cat")
	if category == "" {
		temp.Render(w, "categorie", nil)
		return
	}

	params := c.service.GetCategoryParams(category)
	data, err := c.service.FetchMarvelCharactersWithParams(0, 20, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	characters, err := c.service.ProcessCharacters([]byte(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp.Render(w, "categorie", map[string]interface{}{
		"data": map[string]interface{}{"results": characters},
	})
}

func (c *CharacterController) HandleCharacterDetails(w http.ResponseWriter, r *http.Request) {
	// Implémentation de la gestion des détails d'un personnage
	// ... à compléter selon vos besoins
}
