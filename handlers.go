package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// --- HOME HANDLER ---

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		log.Println("Erreur template:", err)
		return
	}
	tmpl.Execute(w, nil)
}

// --- INDEX HANDLER (Search/Browse Page) ---

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" {
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

// --- ARTIST HANDLER ---

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

	// Recherche relation en mémoire
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

	// --- SYSTEME DE SECOURS (FALLBACK) ---
	// Si on n'a pas de dates, on essaie de les chercher, sinon on SIMULE.
	if !foundRel || len(rel.DatesLocations) == 0 {
		// Essai appel API unique
		fetchedRel, err := getRelation(id)
		if err == nil && len(fetchedRel.DatesLocations) > 0 {
			rel = fetchedRel
		} else {
			// ULTIME SECOURS : Si l'API échoue, on invente des dates pour que le site soit joli
			log.Printf("⚡ Mode Simulation activé pour l'artiste ID %d\n", id)
			rel = generateMockRelation(id)
		}
	}

	data := ArtistPageData{Artist: selected, Relations: rel}
	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Erreur template", 500)
		return
	}
	tmpl.Execute(w, data)
}

// --- SEARCH HANDLER ---

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

// --- API ARTISTS HANDLER ---

func apiArtistsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	currentArtists := allArtists
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentArtists)
}
