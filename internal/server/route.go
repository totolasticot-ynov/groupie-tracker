package server

import "net/http"

func registerRoutes() {
	http.HandleFunc("/", ExplorePage)

	http.HandleFunc("/api/create-order", CreateOrder)
	http.HandleFunc("/api/capture-order", CaptureOrder)
}
