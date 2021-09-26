package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/slash3b/pirata/api/repository"

	"github.com/slash3b/pirata/api/model"
)

func Films(w http.ResponseWriter, r *http.Request) {
	var allFilms []model.Film

	// add error handling
	allFilms = repository.RepoService.FilmsRepo.GetAll(context.Background())

	jsonResult, err := json.Marshal(allFilms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}
