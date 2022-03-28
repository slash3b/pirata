package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/slash3b/pirata/api/http/handlers"
	"github.com/slash3b/pirata/api/http/middleware"

	"github.com/slash3b/pirata/api/repository/repos"

	"github.com/slash3b/pirata/api/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed static/api/index.html
var indexHtml []byte

//go:embed openapi.json
var openapiJson []byte

func main() {
	service, err := initRepoService()
	if err != nil {
		log.Fatal("unable to initialize repositories ", err)
	}

	repository.RepoService = service

	mux := initMux()

	fmt.Println("starting http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}

func initMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/films", middleware.RateLimiter(10, handlers.Films))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHtml)
	})

	mux.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openapiJson)
	})

	return mux
}

func initRepoService() (*repository.Service, error) {
	db, err := gorm.Open(sqlite.Open("../pirata.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	filmsRepo := repos.NewFilmsRepository(db)

	return repository.NewService(filmsRepo), nil
}
