package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"gorm.io/gorm"
)

type ITransferRepository interface {
	GetByID(c *gin.Context, id uint) (models.Transfer, error)
	Save(c *gin.Context, acc models.Transfer) (models.Transfer, error)
	GetPendingPaymentsAmountForAccount(c *gin.Context, accountID string) (uint, error)
	SetStatus(c *gin.Context, id uint, newStatus models.TransferStatus) error
	CompleteInternalTransfer(c *gin.Context, transfer models.Transfer) error
	MakeExternalPayment(c *gin.Context, accountID string, transfer models.Transfer) error
	ReceiveExternalPayment(c *gin.Context, accountID string, transfer models.Transfer) error
}

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) transferRepository {
	return transferRepository{db}
}

var _ ITransferRepository = (transferRepository{})

func (r transferRepository) GetByID(c *gin.Context, id uint) (models.Transfer, error) {
	var item models.Transfer
	if result := r.db.First(&item, id); result.Error != nil {
		log.Error(fmt.Sprintf("transferRepository | GetByID err - %s", result.Error.Error()))
		return models.Transfer{}, result.Error
	}

	return item, nil
}

func (r transferRepository) Save(c *gin.Context, tr models.Transfer) (models.Transfer, error) {
	if result := r.db.Save(&tr); result.Error != nil {
		log.Error(fmt.Sprintf("transferRepository | Save err - %s", result.Error.Error()))
		//todo type this error
		return models.Transfer{}, result.Error
	}

	return tr, nil
}
func (r transferRepository) GetPendingPaymentsAmountForAccount(c *gin.Context, accountID string) (uint, error) {
	var totalPending *uint
	if result := r.db.Model(&models.Transfer{}).
		Select("SUM(amount)").
		Where("from_account_id = ? AND status = ?", accountID, "pending").
		Scan(&totalPending); result.Error != nil {

		log.Error(fmt.Sprintf("transferRepository | GetPendingPaymentsAmountForAccount err - %s", result.Error.Error()))
		return 0, result.Error
	}

	if totalPending == nil {
		return 0, nil
	}

	return uint(*totalPending), nil
}

func (r transferRepository) SetStatus(c *gin.Context, id uint, newStatus models.TransferStatus) error {
	if err := r.db.Model(&models.Transfer{}).
		Where("id = ?", id).
		Update("status", newStatus).Error; err != nil {
		log.Error(fmt.Sprintf("transferRepository | SetStatus err - %s", err.Error()))
		return err
	}

	return nil
}

func (r transferRepository) CompleteInternalTransfer(
	c *gin.Context,
	transfer models.Transfer,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// decrease from acc balance
		//TODO IMPORTANT this shouldnt be done by this repo object
		if err := tx.Model(&models.Account{}).
			Where("id = ?", transfer.FromAccountID).
			Update("balance", gorm.Expr("balance - ?", transfer.Amount)).Error; err != nil {
			log.Error(fmt.Sprintf("CompleteTransfer | Decrease fromAccount failed - %s", err.Error()))
			return err
		}

		// increase to acc balance
		//TODO IMPORTANT this shouldnt be done by this repo object
		if err := tx.Model(&models.Account{}).
			Where("id = ?", transfer.ToAccountID).
			Update("balance", gorm.Expr("balance + ?", transfer.Amount)).Error; err != nil {
			log.Error(fmt.Sprintf("CompleteTransfer | Increase toAccount failed - %s", err.Error()))
			return err
		}

		// update transfer status to "completed"
		if err := tx.Model(&models.Transfer{}).
			Where("id = ?", transfer.ID).
			Update("status", models.TRANSFER_STATUS_COMPLETED).Error; err != nil {
			log.Error(fmt.Sprintf("CompleteTransfer | Update transfer status failed - %s", err.Error()))
			return err
		}

		return nil
	})
}

func (r transferRepository) MakeExternalPayment(
	c *gin.Context,
	accountID string,
	transfer models.Transfer,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Subtract transfer amount from account
		if err := tx.Model(&models.Account{}).
			Where("id = ?", accountID).
			Update("balance", gorm.Expr("balance - ?", transfer.Amount)).Error; err != nil {
			log.Error(fmt.Sprintf("transferRepository | MakePayment - balance update err - %s", err.Error()))
			return err
		}

		// Mark transfer as completed
		if err := tx.Model(&models.Transfer{}).
			Where("id = ?", transfer.ID).
			Update("status", models.TRANSFER_STATUS_COMPLETED).Error; err != nil {
			log.Error(fmt.Sprintf("transferRepository | MakePayment - status update err - %s", err.Error()))
			return err
		}

		return nil // success: commit
	})
}

func (r transferRepository) ReceiveExternalPayment(
	c *gin.Context,
	accountID string,
	transfer models.Transfer,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Subtract transfer amount from account
		if err := tx.Model(&models.Account{}).
			Where("id = ?", accountID).
			Update("balance", gorm.Expr("balance + ?", transfer.Amount)).Error; err != nil {
			log.Error(fmt.Sprintf("transferRepository | ReceivePayment - balance update err - %s", err.Error()))
			return err
		}

		// Mark transfer as completed
		if err := tx.Model(&models.Transfer{}).
			Where("id = ?", transfer.ID).
			Update("status", models.TRANSFER_STATUS_COMPLETED).Error; err != nil {
			log.Error(fmt.Sprintf("transferRepository | ReceivePayment - status update err - %s", err.Error()))
			return err
		}

		return nil // success: commit
	})
}
