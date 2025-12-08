package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// --- STRUCTURES ---

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

// Structure pour l'API "index" des relations (toutes les relations d'un coup)
type RelationsIndex struct {
	Index []Relations `json:"index"`
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

// Structure des résultats de recherche pour le JSON
type SearchResult struct {
	Text     string `json:"text"`     // Le texte affiché (ex: "Phil Collins")
	Type     string `json:"type"`     // Le type (ex: "member")
	ArtistId int    `json:"artistId"` // Pour la redirection
}

// --- VARIABLES GLOBALES (CACHE) ---
var (
	allArtists   []Artist
	allRelations []Relations
	mutex        sync.Mutex
)

// --- MAIN ---

func main() {
	// Chargement initial des données pour la recherche
	log.Println("Chargement des données API...")
	loadData()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/search", searchHandler) // Nouvelle route pour la barre de recherche

	log.Println("✅ Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- LOGIQUE DE CHARGEMENT ---
func loadData() {
	mutex.Lock()
	defer mutex.Unlock()

	// 1. Récupérer les artistes
	apiArtists, _ := getArtists()
	myArtists := getCustomArtists()
	allArtists = append(myArtists, apiArtists...)

	// 2. Récupérer TOUTES les relations (pour la recherche par lieu)
	// Note: Pour tes artistes persos, on simule, pour l'API on fetch tout l'index
	apiRelIndex, err := getAllRelationsIndex()
	if err == nil {
		allRelations = apiRelIndex.Index
	} else {
		log.Println("Erreur chargement relations API:", err)
	}

	// Ajout des relations pour tes artistes persos (simulées pour la recherche)
	for _, art := range myArtists {
		allRelations = append(allRelations, Relations{
			Id: art.Id,
			DatesLocations: map[string][]string{
				"london-uk": {"01-01-2024"}, "paris-france": {"05-01-2024"},
			},
		})
	}
}

// --- HANDLERS ---

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	results := []SearchResult{}

	if query == "" {
		json.NewEncoder(w).Encode(results)
		return
	}

	mutex.Lock() // Lecture sûre des variables globales
	artists := allArtists
	relations := allRelations
	mutex.Unlock()

	// Map pour éviter les doublons exacts
	seen := make(map[string]bool)

	for _, a := range artists {
		// 1. Recherche par NOM (Artist/Band)
		if strings.Contains(strings.ToLower(a.Name), query) {
			addResult(&results, a.Name, "artist/band", a.Id, seen)
		}

		// 2. Recherche par MEMBRE
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), query) {
				addResult(&results, m, "member", a.Id, seen)
			}
		}

		// 3. Recherche par DATE CREATION
		if strings.Contains(strconv.Itoa(a.CreationDate), query) {
			addResult(&results, "Créé en "+strconv.Itoa(a.CreationDate), "creation date", a.Id, seen)
		}

		// 4. Recherche par PREMIER ALBUM
		if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
			addResult(&results, "1er Album: "+a.FirstAlbum, "first album", a.Id, seen)
		}
	}

	// 5. Recherche par LOCATION
	// On doit lier l'ID de la relation à l'artiste
	for _, rel := range relations {
		for loc := range rel.DatesLocations {
			cleanLoc := strings.ReplaceAll(strings.ReplaceAll(loc, "-", " "), "_", " ")
			if strings.Contains(strings.ToLower(cleanLoc), query) {
				// Trouver le nom de l'artiste correspondant à cet ID
				artistName := "Inconnu"
				for _, a := range artists {
					if a.Id == rel.Id {
						artistName = a.Name
						break
					}
				}
				// Affichage: "Paris, France (Queen)"
				text := cleanLoc + " (" + artistName + ")"
				addResult(&results, text, "location", rel.Id, seen)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func addResult(results *[]SearchResult, text, typeStr string, id int, seen map[string]bool) {
	// Clé unique pour éviter les doublons dans l'affichage
	key := text + typeStr + strconv.Itoa(id)
	if !seen[key] {
		*results = append(*results, SearchResult{Text: text, Type: typeStr, ArtistId: id})
		seen[key] = true
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Utilisation des données en cache + filtres dynamiques
	mutex.Lock()
	currentArtists := allArtists
	mutex.Unlock()

	filteredArtists := []Artist{}

	minCreation, _ := strconv.Atoi(r.URL.Query().Get("minCreation"))
	maxCreation, _ := strconv.Atoi(r.URL.Query().Get("maxCreation"))
	minAlbum, _ := strconv.Atoi(r.URL.Query().Get("minAlbum"))
	maxAlbum, _ := strconv.Atoi(r.URL.Query().Get("maxAlbum"))
	members := r.URL.Query()["members"]

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

	for _, a := range currentArtists {
		keep := true
		if a.CreationDate < minCreation || a.CreationDate > maxCreation {
			keep = false
		}
		albYear := extractYear(a.FirstAlbum)
		if albYear < minAlbum || albYear > maxAlbum {
			keep = false
		}
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
		http.Error(w, "Erreur: templates/index.html introuvable", 500)
		return
	}
	tmpl.Execute(w, data)
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	var selected Artist
	mutex.Lock()
	for _, a := range allArtists {
		if a.Id == id {
			selected = a
			break
		}
	}
	// Récupérer la relation spécifique (depuis le cache ou l'API si besoin, ici on simplifie en cherchant dans le cache)
	var rel Relations
	foundRel := false
	for _, r := range allRelations {
		if r.Id == id {
			rel = r
			foundRel = true
			break
		}
	}
	mutex.Unlock()

	if !foundRel {
		// Fallback si pas trouvé dans le cache (ex: appel direct)
		if id < 100 {
			rel, _ = getRelation(id)
		}
	}

	data := ArtistPageData{Artist: selected, Relations: rel}
	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Erreur: templates/artist.html introuvable", 500)
		return
	}
	tmpl.Execute(w, data)
}

// --- DATA & API ---

func extractYear(date string) int {
	parts := strings.Split(date, "-")
	if len(parts) == 3 {
		y, _ := strconv.Atoi(parts[2])
		return y
	}
	return 2000
}

func getAllRelationsIndex() (RelationsIndex, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relations")
	if err != nil {
		return RelationsIndex{}, err
	}
	defer resp.Body.Close()
	var index RelationsIndex
	json.NewDecoder(resp.Body).Decode(&index)
	return index, nil
}

func getCustomArtists() []Artist {
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
