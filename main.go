package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// 1. Structure Artist (Standard) [cite: 39]
type Artist struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// 2. Structure Relations avec la Map demandée [cite: 16, 18, 19]
// Les clés sont dynamiques (ex: "dunedin-new_zealand"), donc on utilise map[string][]string
type Relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Structure globale pour envoyer les données à la page HTML
type PageData struct {
	Artists   []Artist
	Relations Relations // On ajoute une relation pour l'exemple du PDF
}

func main() {
	// Servir les fichiers statiques (CSS/JS) [cite: 51]
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// A. Récupération des Artistes (Simulé ou via API réelle)
	// Pour l'exemple, je récupère la vraie API
	artists, err := getArtists()
	if err != nil {
		http.Error(w, "Erreur API Artistes", 500)
		return
	}

	// B. Récupération d'une Relation (Pour tester la Map du PDF)
	// Je prends la relation de l'ID 1 pour l'exemple
	relation, err := getRelation(1)
	if err != nil {
		http.Error(w, "Erreur API Relations", 500)
		return
	}

	data := PageData{
		Artists:   artists,
		Relations: relation,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl.Execute(w, data)
}

// Fonctions utilitaires pour contacter l'API
func getArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var artists []Artist
	json.NewDecoder(resp.Body).Decode(&artists)
	return artists, nil
}

func getRelation(id int) (Relations, error) {
	// Note: L'API réelle retourne parfois une structure imbriquée,
	// mais ici on adapte selon la structure stricte demandée par ton PDF.
	// URL exemple : https://groupietrackers.herokuapp.com/api/relations/1
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relations/1")
	if err != nil {
		return Relations{}, err
	}
	defer resp.Body.Close()

	var result Relations
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
