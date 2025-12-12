package main

import (
	"strconv"
	"strings"
)

// --- UTILITY FUNCTIONS ---

func extractYear(d string) int {
	parts := strings.Split(d, "-")
	if len(parts) == 3 {
		v, _ := strconv.Atoi(parts[2])
		return v
	}
	return 2000
}

// --- MOCK DATA GENERATOR ---

func generateMockRelation(id int) Relations {
	return Relations{
		Id: id,
		DatesLocations: map[string][]string{
			"los_angeles-usa": {"20-10-2024", "21-10-2024"},
			"paris-france":    {"15-11-2024"},
			"london-uk":       {"01-12-2024"},
			"berlin-germany":  {"05-01-2025"},
		},
	}
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
