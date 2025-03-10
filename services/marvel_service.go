package services

import (
	"api/models"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Constantes pour l'authentification à l'API Marvel
const (
	publicKey  = "1e2027484c06e5cc452a3b8de485b482"
	privateKey = "35eb641e8155bef52b1fd281b32d16611b1841c7"
	baseURL    = "https://gateway.marvel.com/v1/public"

	// Constantes pour les séries spécifiques
	AVENGERS_SERIES = "24229,21366,9085"    // IDs des séries Avengers
	XMEN_SERIES     = "2258,27567,2258"     // IDs des séries X-Men
	POPULAR_SERIES  = "454,354,31445,16452" // IDs des séries populaires
)

type MarvelService struct{}

func NewMarvelService() *MarvelService {
	return &MarvelService{}
}

// generateAuthParams génère les paramètres d'authentification nécessaires pour l'API Marvel
func generateAuthParams() (string, string, string) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	token := ts + privateKey + publicKey
	hash := md5.Sum([]byte(token))
	hashStr := hex.EncodeToString(hash[:])
	return ts, hashStr, publicKey
}

func (s *MarvelService) FetchMarvelCharacters(offset int, limit int) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	url := fmt.Sprintf("%s/characters?ts=%s&apikey=%s&hash=%s&offset=%d&limit=%d&orderBy=name",
		baseURL, ts, pubKey, hashStr, offset, limit)

	return fetchFromAPI(url)
}

func (s *MarvelService) SearchCharacters(query string) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	url := fmt.Sprintf("%s/characters?nameStartsWith=%s&ts=%s&apikey=%s&hash=%s",
		baseURL, query, ts, pubKey, hashStr)

	return fetchFromAPI(url)
}

func (s *MarvelService) FetchCharacterComics(characterID int) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	url := fmt.Sprintf("%s/characters/%d/comics?ts=%s&apikey=%s&hash=%s",
		baseURL, characterID, ts, pubKey, hashStr)

	return fetchFromAPI(url)
}

func (s *MarvelService) FetchCharacterSeries(characterID int) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	url := fmt.Sprintf("%s/characters/%d/series?ts=%s&apikey=%s&hash=%s",
		baseURL, characterID, ts, pubKey, hashStr)

	return fetchFromAPI(url)
}

func (s *MarvelService) GetCategoryParams(category string) string {
	switch category {
	case "avengers":
		return "&series=24229,21366,9085"
	case "xmen":
		return "&series=2258,27567,2258"
	case "other":
		return "&orderBy=-modified&limit=20"
	default:
		return ""
	}
}

func (s *MarvelService) FetchMarvelCharactersWithParams(offset int, limit int, additionalParams string) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()

	baseQuery := fmt.Sprintf("%s/characters?ts=%s&apikey=%s&hash=%s&offset=%d&limit=%d",
		baseURL, ts, pubKey, hashStr, offset, limit)

	cleanParams := strings.ReplaceAll(additionalParams, " ", "%20")
	url := baseQuery + cleanParams

	return fetchFromAPI(url)
}

func (s *MarvelService) FetchCharacterDetails(id int) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	url := fmt.Sprintf("%s/characters/%d?ts=%s&apikey=%s&hash=%s",
		baseURL, id, ts, pubKey, hashStr)

	// Ajouter un délai entre les requêtes pour respecter les limites de l'API
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 409 {
			// Attendre un peu et réessayer
			time.Sleep(1 * time.Second)
			return s.FetchCharacterDetails(id)
		}
		return "", fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// fetchFromAPI effectue une requête HTTP vers l'API Marvel avec gestion des erreurs et retries
func fetchFromAPI(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		// Gestion spéciale pour le rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("API returned status code: %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(body), nil
	}

	return "", lastErr
}

func (s *MarvelService) ProcessCharacters(data []byte) ([]models.Character, error) {
	var response struct {
		Data struct {
			Results []models.Character `json:"results"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Data.Results, nil
}

// FetchComicsByType récupère les comics selon leur type
func (s *MarvelService) FetchComicsByType(comicType string, offset int) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	var url string

	switch comicType {
	case "popular":
		url = fmt.Sprintf("%s/comics?series=%s&ts=%s&apikey=%s&hash=%s&offset=%d&limit=20&orderBy=-focDate",
			baseURL, POPULAR_SERIES, ts, pubKey, hashStr, offset)
	case "avengers":
		url = fmt.Sprintf("%s/comics?series=%s&ts=%s&apikey=%s&hash=%s&offset=%d&limit=20&orderBy=-focDate",
			baseURL, AVENGERS_SERIES, ts, pubKey, hashStr, offset)
	case "xmen":
		url = fmt.Sprintf("%s/comics?series=%s&ts=%s&apikey=%s&hash=%s&offset=%d&limit=20&orderBy=-focDate",
			baseURL, XMEN_SERIES, ts, pubKey, hashStr, offset)
	}

	return fetchFromAPI(url)
}

// FetchCharactersByTeam récupère les personnages par équipe
func (s *MarvelService) FetchCharactersByTeam(team string, offset int) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	var url string

	switch team {
	case "avengers":
		url = fmt.Sprintf("%s/characters?series=%s&ts=%s&apikey=%s&hash=%s&offset=%d&limit=20",
			baseURL, AVENGERS_SERIES, ts, pubKey, hashStr, offset)
	case "xmen":
		url = fmt.Sprintf("%s/characters?series=%s&ts=%s&apikey=%s&hash=%s&offset=%d&limit=20",
			baseURL, XMEN_SERIES, ts, pubKey, hashStr, offset)
	}

	return fetchFromAPI(url)
}

// FetchUniverseCharacters récupère les personnages d'un univers spécifique
func (s *MarvelService) FetchUniverseCharacters(universe string, offset int) (string, error) {
	ts, hashStr, pubKey := generateAuthParams()
	var url string

	switch universe {
	case "616":
		// Personnages de l'univers principal Marvel (Terre-616)
		url = fmt.Sprintf("%s/characters?series=24229&ts=%s&apikey=%s&hash=%s&offset=%d&limit=20",
			baseURL, ts, pubKey, hashStr, offset)
	case "ultimate":
		// Personnages de l'univers Ultimate
		url = fmt.Sprintf("%s/characters?series=2258&ts=%s&apikey=%s&hash=%s&offset=%d&limit=20",
			baseURL, ts, pubKey, hashStr, offset)
	}

	return fetchFromAPI(url)
}
