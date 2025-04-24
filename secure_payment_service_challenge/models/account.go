package models

import "time"

type Account struct {
	ID        string `gorm:"primaryKey"`
	Balance   uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
