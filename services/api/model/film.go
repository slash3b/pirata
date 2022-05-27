package model

type Film struct {
	ID          int `gorm:"primary_key"`
	Title       string
	Lang        string
	Dimension   string
	ReleaseDate string
}
