package main

import (
	"api/routes"
	"api/templates"
)

func main() {
	templates.InitTemplates()
	routes.InitServe()
}
