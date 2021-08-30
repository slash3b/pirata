package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/slash3b/pirata/api/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("../pirata_prod.db"), &gorm.Config{})
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

	fmt.Println("starting localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", mux))

}
