package routes

import (
	"fmt"
	"net/http"
)

func InitServe() {
	// Démarrer le serveur
	fmt.Println("Le serveur est accesible via : http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
