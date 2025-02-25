package main

import (
	"api/routes"
	"api/services"
	"api/templates"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialiser les templates
	templates.InitTemplates()

	// Créer le service de cache
	cacheService := services.NewCacheService()

	// Configurer les routes
	routes.InitRoutes(cacheService)

	// Démarrer le serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// value hash
	/*ts := 19
	public_key := "1e2027484c06e5cc452a3b8de485b482"
	private_key := "35eb641e8155bef52b1fd281b32d16611b1841c7"
	token := fmt.Sprintf("%d%s%s", ts, private_key, public_key)
	// exemple md5
	h := md5.New()
	h.Write([]byte(token))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(hash)*/
}
