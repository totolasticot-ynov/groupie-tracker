package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Mise à jour de la structure pour inclure les données nécessaires aux filtres
type Artist struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	CreationDate int      `json:"creationDate"` // Année de création
	FirstAlbum   string   `json:"firstAlbum"`   // Date premier album (ex: "14-02-1999")
	Members      []string `json:"members"`      // Liste des membres
}

type Relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type PageData struct {
	Artists   []Artist
	Relations Relations
	Filters   FilterData // Pour garder les valeurs du formulaire affichées
}

// Structure pour réafficher ce que l'utilisateur a coché
type FilterData struct {
	MinCreation int
	MaxCreation int
	MinAlbum    int
	MaxAlbum    int
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 1. Récupérer tous les artistes (API + Custom)
	apiArtists, err := getArtists()
	if err != nil {
		log.Println("Erreur API:", err)
	}
	myArtists := getCustomArtists()
	allArtists := append(myArtists, apiArtists...)

	// 2. Gestion des Filtres (Si le formulaire est envoyé)
	filteredArtists := []Artist{}

	// Récupération des paramètres de l'URL
	minCreation, _ := strconv.Atoi(r.URL.Query().Get("minCreation"))
	maxCreation, _ := strconv.Atoi(r.URL.Query().Get("maxCreation"))
	minAlbumYear, _ := strconv.Atoi(r.URL.Query().Get("minAlbum"))
	maxAlbumYear, _ := strconv.Atoi(r.URL.Query().Get("maxAlbum"))
	membersSelected := r.URL.Query()["members"] // Récupère les cases cochées (ex: ["1", "4"])

	// Si c'est le premier chargement (pas de filtres), on met des valeurs par défaut larges
	if r.URL.Query().Get("minCreation") == "" {
		minCreation = 1950
		maxCreation = 2025
		minAlbumYear = 1950
		maxAlbumYear = 2025
	}

	// BOUCLE DE FILTRAGE
	for _, artist := range allArtists {
		keep := true

		// Filtre 1 : Date de création (Range) [cite: 65, 71]
		if artist.CreationDate < minCreation || artist.CreationDate > maxCreation {
			keep = false
		}

		// Filtre 2 : Date premier album (Range) [cite: 66]
		// On extrait l'année de la string "dd-mm-yyyy"
		albumYear := extractYear(artist.FirstAlbum)
		if albumYear < minAlbumYear || albumYear > maxAlbumYear {
			keep = false
		}

		// Filtre 3 : Nombre de membres (Checkbox) [cite: 67, 72]
		if len(membersSelected) > 0 {
			nbMembers := strconv.Itoa(len(artist.Members))
			found := false
			for _, m := range membersSelected {
				if m == nbMembers {
					found = true
					break
				}
			}
			if !found {
				keep = false
			}
		}

		if keep {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Récupération Relation (inchangé)
	relation, _ := getRelation(1)

	data := PageData{
		Artists:   filteredArtists, // On envoie la liste filtrée
		Relations: relation,
		Filters: FilterData{
			MinCreation: minCreation,
			MaxCreation: maxCreation,
			MinAlbum:    minAlbumYear,
			MaxAlbum:    maxAlbumYear,
		},
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.Execute(w, data)
}

// Utilitaire pour extraire l'année "1999" de "12-02-1999"
func extractYear(dateStr string) int {
	parts := strings.Split(dateStr, "-")
	if len(parts) == 3 {
		year, _ := strconv.Atoi(parts[2])
		return year
	}
	// Si pas de date (ex: tes artistes custom sans info), on renvoie une valeur par défaut
	return 2000
}

// Tes artistes personnalisés (J'ai ajouté des fausses dates pour que les filtres marchent sur eux aussi)
func getCustomArtists() []Artist {
	return []Artist{
		{Id: 101, Name: "Aphex Twin", Image: "/static/img/aphex.jpg", CreationDate: 1985, FirstAlbum: "01-01-1992", Members: []string{"Richard D. James"}},
		{Id: 102, Name: "Crystal Castles", Image: "/static/img/crystal.jpg", CreationDate: 2003, FirstAlbum: "18-03-2008", Members: []string{"Ethan Kath", "Alice Glass"}},
		{Id: 103, Name: "Ennio Morricone", Image: "/static/img/ennio.jpg", CreationDate: 1946, FirstAlbum: "01-01-1961", Members: []string{"Ennio"}},
		{Id: 104, Name: "Rihanna", Image: "/static/img/rihanna.jpg", CreationDate: 2003, FirstAlbum: "30-08-2005", Members: []string{"Rihanna"}},
		{Id: 105, Name: "Daft Punk", Image: "/static/img/daftpunk.jpg", CreationDate: 1993, FirstAlbum: "20-01-1997", Members: []string{"Thomas", "Guy-Manuel"}},
		{Id: 106, Name: "TV Girl", Image: "/static/img/tvgirl.jpg", CreationDate: 2010, FirstAlbum: "01-01-2014", Members: []string{"Brad", "Jason", "Wyatt"}},
		{Id: 107, Name: "Björk", Image: "/static/img/bjork.jpg", CreationDate: 1977, FirstAlbum: "05-07-1993", Members: []string{"Björk"}},
		{Id: 108, Name: "Nirvana", Image: "/static/img/nirvana.jpg", CreationDate: 1987, FirstAlbum: "15-06-1989", Members: []string{"Kurt", "Dave", "Krist"}},
		{Id: 109, Name: "Snow Strippers", Image: "/static/img/snow.jpg", CreationDate: 2021, FirstAlbum: "01-01-2022", Members: []string{"Tatiana", "Graham"}},
		{Id: 110, Name: "Venetian Snares", Image: "/static/img/venetian.jpg", CreationDate: 1992, FirstAlbum: "01-01-1998", Members: []string{"Aaron Funk"}},
		{Id: 111, Name: "Boards of Canada", Image: "/static/img/boc.jpg", CreationDate: 1986, FirstAlbum: "01-01-1998", Members: []string{"Mike", "Marcus"}},
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

func getRelation(id int) (Relations, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relations/1")
	if err != nil {
		return Relations{}, err
	}
	defer resp.Body.Close()
	var result Relations
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
