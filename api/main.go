package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/slash3b/pirata/api/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed static/api/index.html
var indexHtml []byte

//go:embed openapi.json
var openapiJson []byte

func main() {
	db, err := gorm.Open(sqlite.Open("../pirata.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("unable to connect to sqlite db", err)
	}

	mux := http.NewServeMux()

	rateLimiterMiddleware := func(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

		limitingBuffer := make(chan struct{}, 10)
		return func(w http.ResponseWriter, r *http.Request) {

			limitingBuffer <- struct{}{}
			defer func() { <-limitingBuffer }()

			f(w, r)
		}
	}

	getFilms := func(w http.ResponseWriter, r *http.Request) {
		var allFilms []model.Film

		result := db.Order("register_date desc").Find(&allFilms)
		if result.Error != nil {
			panic("ooops")
		}

		fmt.Println("rows affected ", result.RowsAffected)

		jsonResult, err := json.Marshal(allFilms)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)
	}

	mux.HandleFunc("/films", rateLimiterMiddleware(getFilms))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHtml)
	})

	mux.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openapiJson)
	})

	fmt.Println("starting http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
