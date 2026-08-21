package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jb843051627/tephra-chron/internal/handler"
	"github.com/jb843051627/tephra-chron/internal/service"
	"github.com/jb843051627/tephra-chron/internal/store"
)

func main() {
	path := os.Getenv("TEPHRA_CHRON_DB")
	if path == "" {
		path = "data/tephra-chron.db"
	}
	repository, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	app := service.NewLab(repository)
	defer app.Close()
	addr := os.Getenv("TEPHRA_CHRON_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("tephra-chron listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.New(app)))
}
