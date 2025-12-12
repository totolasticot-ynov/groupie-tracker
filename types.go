package main

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
