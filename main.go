package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// --- STRUCTURES DE DONNÉES ---

type Artist struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Members      []string `json:"members"`
}

type Relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type PageData struct {
	Artists []Artist
	Filters FilterData
}

type FilterData struct {
	MinCreation int
	MaxCreation int
	MinAlbum    int
	MaxAlbum    int
}

type ArtistPageData struct {
	Artist    Artist
	Relations Relations
}

// --- MAIN ---

func main() {
	// 1. Servir le dossier "static" (pour les images et scripts)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Les Routes
	http.HandleFunc("/", indexHandler)        // Page d'accueil (Liste + Filtres)
	http.HandleFunc("/artist", artistHandler) // Page Détail (Carte + Infos)

	log.Println("✅ Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- HANDLERS (LOGIQUE) ---

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Récupérer les artistes (API + Custom)
	apiArtists, _ := getArtists()
	myArtists := getCustomArtists()
	allArtists := append(myArtists, apiArtists...)

	// -- LOGIQUE DE FILTRAGE --
	filteredArtists := []Artist{}

	// Récupération des paramètres URL
	minCreation, _ := strconv.Atoi(r.URL.Query().Get("minCreation"))
	maxCreation, _ := strconv.Atoi(r.URL.Query().Get("maxCreation"))
	minAlbum, _ := strconv.Atoi(r.URL.Query().Get("minAlbum"))
	maxAlbum, _ := strconv.Atoi(r.URL.Query().Get("maxAlbum"))
	members := r.URL.Query()["members"]

	// Valeurs par défaut si le formulaire n'est pas utilisé
	if minCreation == 0 {
		minCreation = 1950
	}
	if maxCreation == 0 {
		maxCreation = 2025
	}
	if minAlbum == 0 {
		minAlbum = 1950
	}
	if maxAlbum == 0 {
		maxAlbum = 2025
	}

	// Boucle de filtrage
	for _, a := range allArtists {
		keep := true

		// 1. Filtre Date Création
		if a.CreationDate < minCreation || a.CreationDate > maxCreation {
			keep = false
		}

		// 2. Filtre Année Album
		albYear := extractYear(a.FirstAlbum)
		if albYear < minAlbum || albYear > maxAlbum {
			keep = false
		}

		// 3. Filtre Nombre de Membres
		if len(members) > 0 {
			nbStr := strconv.Itoa(len(a.Members))
			found := false
			for _, m := range members {
				if m == nbStr {
					found = true
					break
				}
			}
			if !found {
				keep = false
			}
		}

		if keep {
			filteredArtists = append(filteredArtists, a)
		}
	}

	data := PageData{
		Artists: filteredArtists,
		Filters: FilterData{MinCreation: minCreation, MaxCreation: maxCreation, MinAlbum: minAlbum, MaxAlbum: maxAlbum},
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur critique : templates/index.html introuvable", 500)
		return
	}
	tmpl.Execute(w, data)
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	// Retrouver l'artiste cliqué
	var selected Artist
	apiArtists, _ := getArtists()
	myArtists := getCustomArtists()
	allArtists := append(myArtists, apiArtists...)

	for _, a := range allArtists {
		if a.Id == id {
			selected = a
			break
		}
	}

	// Retrouver ses relations (dates/lieux)
	var rel Relations
	if id < 100 {
		// Appel API réel
		rel, _ = getRelation(id)
	} else {
		// Données simulées pour tes artistes persos
		rel = Relations{
			Id: id,
			DatesLocations: map[string][]string{
				"london-uk":    {"01-01-2024"},
				"paris-france": {"05-01-2024", "06-01-2024"},
				"tokyo-japan":  {"10-01-2024"},
			},
		}
	}

	data := ArtistPageData{
		Artist:    selected,
		Relations: rel,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Erreur critique : templates/artist.html introuvable", 500)
		return
	}
	tmpl.Execute(w, data)
}

// --- FONCTIONS UTILITAIRES ---

func extractYear(date string) int {
	parts := strings.Split(date, "-")
	if len(parts) == 3 {
		y, _ := strconv.Atoi(parts[2])
		return y
	}
	return 2000
}

func getCustomArtists() []Artist {
	// Liste de tes artistes personnalisés
	return []Artist{
		{Id: 101, Name: "Aphex Twin", Image: "/static/img/aphex.jpg", CreationDate: 1985, FirstAlbum: "01-01-1992", Members: []string{"Richard"}},
		{Id: 102, Name: "Crystal Castles", Image: "/static/img/crystal.jpg", CreationDate: 2003, FirstAlbum: "18-03-2008", Members: []string{"Ethan", "Alice"}},
		{Id: 103, Name: "Ennio Morricone", Image: "/static/img/ennio.jpg", CreationDate: 1946, FirstAlbum: "01-01-1961", Members: []string{"Ennio"}},
		{Id: 104, Name: "Rihanna", Image: "/static/img/rihanna.jpg", CreationDate: 2003, FirstAlbum: "30-08-2005", Members: []string{"Rihanna"}},
		{Id: 105, Name: "Daft Punk", Image: "/static/img/daftpunk.jpg", CreationDate: 1993, FirstAlbum: "20-01-1997", Members: []string{"Thomas", "Guy-Manuel"}},
		{Id: 106, Name: "TV Girl", Image: "/static/img/tvgirl.jpg", CreationDate: 2010, FirstAlbum: "01-01-2014", Members: []string{"Brad", "Jason", "Wyatt"}},
		{Id: 107, Name: "Björk", Image: "/static/img/bjork.jpg", CreationDate: 1977, FirstAlbum: "05-07-1993", Members: []string{"Björk"}},
		{Id: 108, Name: "Snow Strippers", Image: "/static/img/snow.jpg", CreationDate: 2021, FirstAlbum: "01-01-2022", Members: []string{"Tatiana", "Graham"}},
		{Id: 109, Name: "Venetian Snares", Image: "/static/img/venetian.jpg", CreationDate: 1992, FirstAlbum: "01-01-1998", Members: []string{"Aaron"}},
		{Id: 110, Name: "Boards of Canada", Image: "/static/img/boc.jpg", CreationDate: 1986, FirstAlbum: "01-01-1998", Members: []string{"Mike", "Marcus"}},
		{Id: 111, Name: "Deftones", Image: "/static/img/deftones.jpg", CreationDate: 1988, FirstAlbum: "03-10-1995", Members: []string{"Chino", "Stephen", "Abe", "Frank"}},
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
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relations/" + strconv.Itoa(id))
	if err != nil {
		return Relations{}, err
	}
	defer resp.Body.Close()
	var res Relations
	json.NewDecoder(resp.Body).Decode(&res)
	return res, nil
}
