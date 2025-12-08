package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// Structure Artist
type Artist struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// Structure Relations
type Relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type PageData struct {
	Artists   []Artist
	Relations Relations
}

func main() {
	// C'est cette ligne qui permet d'afficher tes images locales
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 1. Récupérer les artistes de l'API officielle
	apiArtists, err := getArtists()
	if err != nil {
		log.Println("Attention : Impossible de contacter l'API (Internet coupé ?)")
		// On continue quand même pour afficher tes artistes perso
	}

	// 2. Récupérer tes artistes personnalisés (Images locales)
	myArtists := getCustomArtists()

	// 3. Fusionner les deux listes (Tes artistes en premier)
	allArtists := append(myArtists, apiArtists...)

	// Récupérer une relation pour l'exemple (ID 1)
	relation, _ := getRelation(1)

	data := PageData{
		Artists:   allArtists,
		Relations: relation,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur : Impossible de lire index.html. Vérifie le dossier templates.", 500)
		return
	}

	tmpl.Execute(w, data)
}

// C'est ici qu'on fait le lien avec tes fichiers dans le dossier "img"
func getCustomArtists() []Artist {
	return []Artist{
		{Id: 101, Name: "Aphex Twin", Image: "/static/img/aphex.jpg"},
		{Id: 102, Name: "Crystal Castles", Image: "/static/img/crystal.jpg"},
		{Id: 103, Name: "Ennio Morricone", Image: "/static/img/ennio.jpg"},
		{Id: 104, Name: "Rihanna", Image: "/static/img/rihanna.jpg"},
		{Id: 105, Name: "Daft Punk", Image: "/static/img/daftpunk.jpg"},
		{Id: 106, Name: "TV Girl", Image: "/static/img/tvgirl.jpg"},
		{Id: 107, Name: "Björk", Image: "/static/img/bjork.jpg"},
		{Id: 108, Name: "Deftones", Image: "/static/img/deftones.jpg"},
		{Id: 109, Name: "Snow Strippers", Image: "/static/img/snow.jpg"},
		{Id: 110, Name: "Venetian Snares", Image: "/static/img/venetian.jpg"},
		{Id: 111, Name: "Boards of Canada", Image: "/static/img/boc.jpg"},
	}
}

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

func getRelation(int) (Relations, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relations/1")
	if err != nil {
		return Relations{}, err
	}
	defer resp.Body.Close()
	var result Relations
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
