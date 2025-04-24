package controller

import "github.com/maty546/secure_payment_service_challenge/models"

type HandleTransferStartRequest struct {
	FromAccountID string `json:"from"`
	ToAccountID   string `json:"to"`
	Amount        uint   `json:"amount"`
}

func (req HandleTransferStartRequest) parseIntoTransferModel() models.Transfer {
	return models.Transfer{FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount:      req.Amount}
}

type HandleTransferResultCallbackRequest struct {
	TransferID uint                  `json:"transfer_id"`
	Status     models.TransferStatus `json:"status"`
}
