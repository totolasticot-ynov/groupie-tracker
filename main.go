package main

import (
	"groupie-tracker/internal/server"
	"log"
)

func main() {
	log.Println("Démarrage minimal de main.go")
	server.Run()
}
