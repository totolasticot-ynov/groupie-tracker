package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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

type SearchResult struct {
	Text     string `json:"text"`
	Type     string `json:"type"`
	ArtistId int    `json:"artistId"`
}

// --- VARIABLES GLOBALES ---
var (
	allArtists   []Artist
	allRelations []Relations
	mutex        sync.Mutex
	// Client HTTP avec timeout pour éviter que le site ne gèle si l'API plante
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

// --- MAIN ---

func main() {
	log.Println("🚀 Démarrage du serveur...")

	// Chargement en tâche de fond pour que le serveur démarre vite
	go loadData()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/search", searchHandler)

	log.Println("✅ Serveur prêt sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- LOGIQUE CHARGEMENT ---

func loadData() {
	log.Println("🔄 Téléchargement des données API...")
	mutex.Lock()
	defer mutex.Unlock()

	// 1. Artistes API
	apiArtists, err := getArtists()
	if err != nil {
		log.Println("❌ Erreur API Artistes:", err)
	}

	// 2. Artistes Custom
	myArtists := getCustomArtists()
	allArtists = append(myArtists, apiArtists...)

	// 3. Relations API (Tentative de tout charger d'un coup)
	apiRelIndex, err := getAllRelationsIndex()
	if err == nil {
		allRelations = apiRelIndex.Index
		log.Printf("✅ %d relations chargées.\n", len(allRelations))
	} else {
		log.Println("⚠️ API Relations lente : chargement au cas par cas activé.")
	}

	// 4. Relations Custom (Simulation dates)
	for _, art := range myArtists {
		allRelations = append(allRelations, Relations{
			Id: art.Id,
			DatesLocations: map[string][]string{
				"london-uk":    {"01-01-2024", "02-01-2024"},
				"paris-france": {"05-06-2024"},
				"new_york-usa": {"10-09-2024"},
				"tokyo-japan":  {"15-12-2024"},
			},
		})
	}
}

// --- HANDLERS ---

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

	// Chercher relations dans le cache mémoire
	var rel Relations
	foundRel := false
	for _, rl := range allRelations {
		if rl.Id == id {
			rel = rl
			foundRel = true
			break
		}
	}
	mutex.Unlock()

	// SYSTEME DE SECOURS (FALLBACK)
	// Si on n'a pas trouvé les dates en mémoire et que c'est un artiste API (ID < 100)
	if (!foundRel || len(rel.DatesLocations) == 0) && id < 100 {
		log.Printf("🔄 Récupération forcée pour ID %d...\n", id)
		fetchedRel, err := getRelation(id)
		if err == nil {
			rel = fetchedRel
		}
	}

	data := ArtistPageData{Artist: selected, Relations: rel}
	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Erreur template artist", 500)
		return
	}
	tmpl.Execute(w, data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	results := []SearchResult{}
	if query == "" {
		json.NewEncoder(w).Encode(results)
		return
	}

	mutex.Lock()
	artists := allArtists
	relations := allRelations
	mutex.Unlock()

	seen := make(map[string]bool)

	for _, a := range artists {
		if strings.Contains(strings.ToLower(a.Name), query) {
			addResult(&results, a.Name, "Artiste", a.Id, seen)
		}
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), query) {
				addResult(&results, m, "Membre", a.Id, seen)
			}
		}
		if strings.Contains(strconv.Itoa(a.CreationDate), query) {
			addResult(&results, "Créé en "+strconv.Itoa(a.CreationDate), "Date", a.Id, seen)
		}
		if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
			addResult(&results, "Album: "+a.FirstAlbum, "Album", a.Id, seen)
		}
	}

	// Recherche par lieu
	for _, rel := range relations {
		for loc := range rel.DatesLocations {
			cleanLoc := strings.ReplaceAll(strings.ReplaceAll(loc, "-", " "), "_", " ")
			if strings.Contains(strings.ToLower(cleanLoc), query) {
				artName := "Artiste"
				for _, a := range artists {
					if a.Id == rel.Id {
						artName = a.Name
						break
					}
				}
				addResult(&results, cleanLoc+" ("+artName+")", "Lieu", rel.Id, seen)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func addResult(res *[]SearchResult, txt, typ string, id int, seen map[string]bool) {
	key := txt + typ + strconv.Itoa(id)
	if !seen[key] {
		*res = append(*res, SearchResult{Text: txt, Type: typ, ArtistId: id})
		seen[key] = true
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	mutex.Lock()
	currentArtists := allArtists
	mutex.Unlock()

	filtered := []Artist{}
	minC, _ := strconv.Atoi(r.URL.Query().Get("minCreation"))
	maxC, _ := strconv.Atoi(r.URL.Query().Get("maxCreation"))
	minA, _ := strconv.Atoi(r.URL.Query().Get("minAlbum"))
	maxA, _ := strconv.Atoi(r.URL.Query().Get("maxAlbum"))
	members := r.URL.Query()["members"]

	if minC == 0 {
		minC = 1950
	}
	if maxC == 0 {
		maxC = 2025
	}
	if minA == 0 {
		minA = 1950
	}
	if maxA == 0 {
		maxA = 2025
	}

	for _, a := range currentArtists {
		keep := true
		if a.CreationDate < minC || a.CreationDate > maxC {
			keep = false
		}
		y := extractYear(a.FirstAlbum)
		if y < minA || y > maxA {
			keep = false
		}
		if len(members) > 0 {
			nb := strconv.Itoa(len(a.Members))
			found := false
			for _, m := range members {
				if m == nb {
					found = true
					break
				}
			}
			if !found {
				keep = false
			}
		}
		if keep {
			filtered = append(filtered, a)
		}
	}

	data := PageData{
		Artists: filtered,
		Filters: FilterData{MinCreation: minC, MaxCreation: maxC, MinAlbum: minA, MaxAlbum: maxA},
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur index.html", 500)
		return
	}
	tmpl.Execute(w, data)
}

// --- APPELS API & DATA ---

func extractYear(d string) int {
	parts := strings.Split(d, "-")
	if len(parts) == 3 {
		v, _ := strconv.Atoi(parts[2])
		return v
	}
	return 2000
}

func getAllRelationsIndex() (RelationsIndex, error) {
	resp, err := httpClient.Get("https://groupietrackers.herokuapp.com/api/relations")
	if err != nil {
		return RelationsIndex{}, err
	}
	defer resp.Body.Close()
	var idx RelationsIndex
	json.NewDecoder(resp.Body).Decode(&idx)
	return idx, nil
}

func getArtists() ([]Artist, error) {
	resp, err := httpClient.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var arr []Artist
	json.NewDecoder(resp.Body).Decode(&arr)
	return arr, nil
}

// Retry logic pour les relations individuelles
func getRelation(id int) (Relations, error) {
	url := "https://groupietrackers.herokuapp.com/api/relations/" + strconv.Itoa(id)
	for i := 0; i < 2; i++ { // 2 essais
		resp, err := httpClient.Get(url)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			var r Relations
			if err := json.NewDecoder(resp.Body).Decode(&r); err == nil {
				return r, nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return Relations{}, http.ErrHandlerTimeout
}

func getCustomArtists() []Artist {
	return []Artist{
		{Id: 101, Name: "Aphex Twin", Image: "/static/img/aphex.jpg", CreationDate: 1985, FirstAlbum: "01-01-1992", Members: []string{"Richard D. James"}},
		{Id: 102, Name: "Crystal Castles", Image: "/static/img/crystal.jpg", CreationDate: 2003, FirstAlbum: "18-03-2008", Members: []string{"Ethan Kath", "Alice Glass"}},
		{Id: 103, Name: "Ennio Morricone", Image: "/static/img/ennio.jpg", CreationDate: 1946, FirstAlbum: "01-01-1961", Members: []string{"Ennio Morricone"}},
		{Id: 104, Name: "Rihanna", Image: "/static/img/rihanna.jpg", CreationDate: 2003, FirstAlbum: "30-08-2005", Members: []string{"Rihanna"}},
		{Id: 105, Name: "Daft Punk", Image: "/static/img/daftpunk.jpg", CreationDate: 1993, FirstAlbum: "20-01-1997", Members: []string{"Thomas Bangalter", "Guy-Manuel"}},
		{Id: 106, Name: "TV Girl", Image: "/static/img/tvgirl.jpg", CreationDate: 2010, FirstAlbum: "01-01-2014", Members: []string{"Brad Petering", "Jason Wyman"}},
		{Id: 107, Name: "Björk", Image: "/static/img/bjork.jpg", CreationDate: 1977, FirstAlbum: "05-07-1993", Members: []string{"Björk"}},
		{Id: 108, Name: "Snow Strippers", Image: "/static/img/snow.jpg", CreationDate: 2021, FirstAlbum: "01-01-2022", Members: []string{"Tatiana Schwaninger", "Graham Perez"}},
		{Id: 109, Name: "Venetian Snares", Image: "/static/img/venetian.jpg", CreationDate: 1992, FirstAlbum: "01-01-1998", Members: []string{"Aaron Funk"}},
		{Id: 110, Name: "Boards of Canada", Image: "/static/img/boc.jpg", CreationDate: 1986, FirstAlbum: "01-01-1998", Members: []string{"Mike Sandison", "Marcus Sandison"}},
		{Id: 111, Name: "Deftones", Image: "/static/img/deftones.jpg", CreationDate: 1988, FirstAlbum: "03-10-1995", Members: []string{"Chino Moreno", "Stephen Carpenter"}},
	}
}
