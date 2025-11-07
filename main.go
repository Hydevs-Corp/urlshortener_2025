package main

import (
	"log"

	"github.com/axellelanca/urlshortener/cmd"
	_ "github.com/axellelanca/urlshortener/cmd/cli"    // Importe le package 'cli' pour que ses init() soient exécutés
	_ "github.com/axellelanca/urlshortener/cmd/server" // Importe le package 'server' pour que ses init() soient exécutés
)

func main() {
	// Appelle le point d'entrée Cobra défini dans `cmd`.
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %v", r)
		}
	}()

	cmd.Execute()
}
