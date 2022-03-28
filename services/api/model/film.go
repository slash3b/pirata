package model

type Film struct {
	ID           int `gorm:"primary_key"`
	Title        string
	Meta         string
	RegisterDate string // todo: https://stackoverflow.com/questions/42037562/golang-gorm-time-data-type-conversion
}
