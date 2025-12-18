package server

import (
	"html/template"
	"net/http"
)

type PageData struct {
	PayPalClientID string
}

func ExplorePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/explore.html"))

	data := PageData{
		PayPalClientID: "Aao_IAK9WsSbSKqMd-HfOea_SwHvbJAaeJjpXC8eOmwNm5sj6s6kOLUoRSxOaTsnhR8Dr7oflFu2hj4e",
	}

	tmpl.Execute(w, data)
}
