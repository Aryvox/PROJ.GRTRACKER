package routes

import (
	"fmt"
	"net/http"
)

func InitServe() {
	// Récuperation des routes
	Routes()

	fmt.Println("Le serveur est accesible via : http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
