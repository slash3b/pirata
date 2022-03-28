package model

import "time"

type Film struct {
	ID        int
	Title     string
	Lang      string
	Dimension string
	StartDate time.Time
}
