package models

import "time"

type Transfer struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID uint
	ToAccountID   uint
	Amount        int64
	Status        TransferStatus
	ExternalID    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TransferStatus string

var TRANSFER_STATUS_PENDING TransferStatus = "PENDING"
