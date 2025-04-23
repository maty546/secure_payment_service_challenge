package controller

import "github.com/maty546/secure_payment_service_challenge/models"

type HandleTransferStartRequest struct {
	FromAccountID uint  `json:"from"`
	ToAccountID   uint  `json:"to"`
	Amount        int64 `json:"amount"`
}

func (req HandleTransferStartRequest) parseIntoTransferModel() models.Transfer {
	return models.Transfer{FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount:      req.Amount}
}
