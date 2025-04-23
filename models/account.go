package models

import "time"

type Account struct {
	ID        string `gorm:"primaryKey"`
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
