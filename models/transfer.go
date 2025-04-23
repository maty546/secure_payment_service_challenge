package models

import "time"

type Transfer struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID string
	ToAccountID   string
	Amount        uint
	Status        TransferStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TransferStatus string

var TRANSFER_STATUS_PENDING TransferStatus = "PENDING"
