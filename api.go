package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// --- API CALLS ---

func getAllRelationsIndex() (RelationsIndex, error) {
	resp, err := httpClient.Get("https://groupietrackers.herokuapp.com/api/relations")
	if err != nil {
		return RelationsIndex{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RelationsIndex{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var idx RelationsIndex
	if decErr := json.NewDecoder(resp.Body).Decode(&idx); decErr != nil {
		return RelationsIndex{}, decErr
	}
	return idx, nil
}

func getArtists() ([]Artist, error) {
	resp, err := httpClient.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var arr []Artist
	if decErr := json.NewDecoder(resp.Body).Decode(&arr); decErr != nil {
		return nil, decErr
	}
	return arr, nil
}

func getRelation(id int) (Relations, error) {
	resp, err := httpClient.Get("https://groupietrackers.herokuapp.com/api/relations/" + strconv.Itoa(id))
	if err != nil {
		return Relations{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Relations{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var r Relations
	if decErr := json.NewDecoder(resp.Body).Decode(&r); decErr != nil {
		return Relations{}, decErr
	}
	return r, nil
}
