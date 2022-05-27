package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/slash3b/pirata/api/repository/repos"

	"github.com/slash3b/pirata/api/model"
)

var repo *repos.FilmsRepository

func SetRepo(r *repos.FilmsRepository) {
	repo = r
}

func Films(w http.ResponseWriter, r *http.Request) {
	var allFilms []model.Film

	// add error handling
	allFilms, _ = repo.GetAll(context.Background())

	jsonResult, err := json.Marshal(allFilms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}
