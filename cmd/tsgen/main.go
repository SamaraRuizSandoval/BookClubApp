package main

import (
	"log"

	"github.com/SamaraRuizSandoval/BookClubApp/internal/store"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	t := typescriptify.New()

	t.CreateInterface = true
	t.BackupDir = ""

	t.Add(store.Book{})

	if err := t.ConvertToFile("book-club-app/src/models.ts"); err != nil {
		log.Fatal(err)
	}
}
