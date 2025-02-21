package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	publicKey  = "1e2027484c06e5cc452a3b8de485b482"         // Remplace par ta clé publique
	privateKey = "35eb641e8155bef52b1fd281b32d16611b1841c7" // Remplace par ta clé privée
)

func FetchMarvelCharacters() (string, error) {
	// Générer le timestamp
	ts := fmt.Sprintf("%d", time.Now().Unix())

	// Générer le hash
	token := ts + privateKey + publicKey
	hash := md5.Sum([]byte(token))
	hashStr := hex.EncodeToString(hash[:])

	// Construire l'URL de l'API
	url := fmt.Sprintf("https://gateway.marvel.com/v1/public/characters?ts=%s&apikey=%s&hash=%s", ts, publicKey, hashStr)

	// Effectuer la requête HTTP
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Lire la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
