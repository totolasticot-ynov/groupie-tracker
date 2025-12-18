package server

import (
	"html/template"
	"net/http"
	"os"
)

type PageData struct {
	PayPalClientID string
}

func ExplorePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/explore.html"))

	data := PageData{
		PayPalClientID: os.Getenv("PAYPAL_CLIENT_ID"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
	}
}
