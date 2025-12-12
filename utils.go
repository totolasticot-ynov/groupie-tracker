package main

import (
	"strconv"
)

func addResult(res *[]SearchResult, txt, typ string, id int, seen map[string]bool) {
	key := txt + typ + strconv.Itoa(id)
	if !seen[key] {
		*res = append(*res, SearchResult{Text: txt, Type: typ, ArtistId: id})
		seen[key] = true
	}
}
