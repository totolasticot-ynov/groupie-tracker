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
	Image string `json:"image"` // Servira ici pour la pochette de l'album
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
	// Important : On gère le sous-dossier img pour les images
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 1. Récupérer les artistes de l'API
	apiArtists, err := getArtists()
	if err != nil {
		log.Println("Erreur API:", err)
		// On continue même si l'API plante, pour afficher tes artistes perso
	}

	// 2. Créer tes artistes personnalisés
	myArtists := getCustomArtists()

	// 3. Fusionner les deux listes (Les tiens en premier !)
	// On ajoute 'apiArtists' à la suite de 'myArtists'
	allArtists := append(myArtists, apiArtists...)

	// Pour l'exemple des relations (inchangé)
	relation, _ := getRelation(1)

	data := PageData{
		Artists:   allArtists,
		Relations: relation,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl.Execute(w, data)
}

// Ta liste personnalisée avec les pochettes d'albums
func getCustomArtists() []Artist {
	return []Artist{
		{Id: 101, Name: "Aphex Twin", Image: "https://upload.wikimedia.org/wikipedia/en/8/82/Selected_Ambient_Works_85-92.png"},                                      // Selected Ambient Works
		{Id: 102, Name: "Crystal Castles", Image: "https://upload.wikimedia.org/wikipedia/en/a/ae/Crystal_Castles_%28album%29.png"},                                  // Crystal Castles (I)
		{Id: 103, Name: "Ennio Morricone", Image: "https://upload.wikimedia.org/wikipedia/en/0/03/The_Good_the_Bad_and_the_Ugly_soundtrack_cover.jpg"},               // The Good, The Bad...
		{Id: 104, Name: "Rihanna", Image: "https://upload.wikimedia.org/wikipedia/en/3/32/Rihanna_-_Anti.png"},                                                       // Anti
		{Id: 105, Name: "Daft Punk", Image: "https://upload.wikimedia.org/wikipedia/en/a/ae/Daft_Punk_-_Discovery.jpg"},                                              // Discovery
		{Id: 106, Name: "TV Girl", Image: "https://upload.wikimedia.org/wikipedia/en/3/30/TV_Girl_-_French_Exit.png"},                                                // French Exit
		{Id: 107, Name: "Björk", Image: "https://upload.wikimedia.org/wikipedia/en/a/a6/Bjork_Homogenic.png"},                                                        // Homogenic
		{Id: 108, Name: "Nirvana", Image: "https://upload.wikimedia.org/wikipedia/en/b/b7/NirvanaNevermindalbumcover.jpg"},                                           // Nevermind
		{Id: 109, Name: "Snow Strippers", Image: "https://e.snmc.io/i/600/s/5712f55977a4505030245a4911005f77/10839843/snow-strippers-april-mixtape-3-Cover-Art.jpg"}, // April Mixtape 3 (Exemple)
		{Id: 110, Name: "Venetian Snares", Image: "https://upload.wikimedia.org/wikipedia/en/2/22/Rossz_csillag_alatt_született.jpg"},                                // Rossz Csillag...
		{Id: 111, Name: "Boards of Canada", Image: "https://upload.wikimedia.org/wikipedia/en/e/e6/Boards_of_Canada_-_Music_Has_the_Right_to_Children.png"},          // Music Has the Right to Children
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
