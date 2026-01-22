package server

import (
	"html/template"
	"net/http"
)

type ExplorePageData struct {
	PayPalClientID string
}

func ExplorePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/explore.html"))

	data := ExplorePageData{
		PayPalClientID: "Aao_IAK9WsSbSKqMd-HfOea_SwHvbJAaeJjpXC8eOmwNm5sj6s6kOLUoRSxOaTsnhR8Dr7oflFu2hj4e",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
	}
}
