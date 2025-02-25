package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type Favorite struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
	} `json:"thumbnail"`
}

type FavoriteService struct {
	favorites map[int]Favorite
	filePath  string
	mutex     sync.RWMutex
}

func NewFavoriteService() *FavoriteService {
	fs := &FavoriteService{
		favorites: make(map[int]Favorite),
		filePath:  "data/favorites.json",
	}
	fs.loadFavorites()
	return fs
}

func (s *FavoriteService) loadFavorites() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Créer le dossier data s'il n'existe pas
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	// Si le fichier n'existe pas, créer un fichier vide avec un tableau vide
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(s.filePath, []byte("[]"), 0644); err != nil {
			return err
		}
	}

	// Lire le fichier
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	// Désérialiser les données
	var favorites []Favorite
	if err := json.Unmarshal(data, &favorites); err != nil {
		return err
	}

	// Remplir la map
	s.favorites = make(map[int]Favorite)
	for _, fav := range favorites {
		s.favorites[fav.ID] = fav
	}

	return nil
}

func (s *FavoriteService) saveFavorites() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Convertir la map en slice pour la sauvegarde
	favorites := make([]Favorite, 0, len(s.favorites))
	for _, fav := range s.favorites {
		favorites = append(favorites, fav)
	}

	// Sérialiser les données
	data, err := json.MarshalIndent(favorites, "", "    ")
	if err != nil {
		return err
	}

	// Sauvegarder dans le fichier
	return ioutil.WriteFile(s.filePath, data, 0644)
}

func (s *FavoriteService) Add(favorite Favorite) error {
	s.mutex.Lock()
	s.favorites[favorite.ID] = favorite
	s.mutex.Unlock()
	return s.saveFavorites()
}

func (s *FavoriteService) Remove(id int) error {
	s.mutex.Lock()
	delete(s.favorites, id)
	s.mutex.Unlock()
	return s.saveFavorites()
}

func (s *FavoriteService) GetAll() ([]Favorite, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	favorites := make([]Favorite, 0, len(s.favorites))
	for _, fav := range s.favorites {
		favorites = append(favorites, fav)
	}
	return favorites, nil
}

func (s *FavoriteService) IsFavorite(id int) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, exists := s.favorites[id]
	return exists
}

func (s *FavoriteService) ConvertToFavorite(data map[string]interface{}) (*Favorite, error) {
	id := int(data["id"].(float64))
	thumbnail := data["thumbnail"].(map[string]interface{})

	favorite := &Favorite{
		ID:          id,
		Name:        data["name"].(string),
		Description: "",
	}

	if desc, ok := data["description"].(string); ok {
		favorite.Description = desc
	}

	favorite.Thumbnail.Path = thumbnail["path"].(string)
	favorite.Thumbnail.Extension = thumbnail["extension"].(string)

	return favorite, nil
}
