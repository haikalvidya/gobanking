package models

type Currency struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
	Code string `json:"symbol" gorm:"not null"`
}
