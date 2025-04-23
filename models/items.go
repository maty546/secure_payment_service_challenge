package models

type Item struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
