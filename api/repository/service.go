package repository

import "github.com/slash3b/pirata/api/repository/repos"

var RepoService *Service

type Service struct {
	FilmsRepo *repos.FilmsRepository
}

func NewService(f *repos.FilmsRepository) *Service {
	return &Service{FilmsRepo: f}
}
