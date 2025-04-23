package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"gorm.io/gorm"
)

type ITransferRepository interface {
	GetByID(c *gin.Context, id int) (models.Transfer, error)
	Save(c *gin.Context, acc models.Transfer) (models.Transfer, error)
}

type transferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) transferRepository {
	return transferRepository{db}
}

var _ ITransferRepository = (transferRepository{})

func (r transferRepository) GetByID(c *gin.Context, id int) (models.Transfer, error) {
	var item models.Transfer
	if result := r.db.First(&item, id); result.Error != nil {
		log.Error(fmt.Sprintf("transferRepository | GetByID err - %s", result.Error.Error()))
		return models.Transfer{}, result.Error
	}

	return item, nil
}

func (r transferRepository) Save(c *gin.Context, tr models.Transfer) (models.Transfer, error) {
	if result := r.db.Save(&tr); result.Error != nil {
		log.Error(fmt.Sprintf("accountRepository | Save err - %s", result.Error.Error()))
		//todo type this error
		return models.Transfer{}, result.Error
	}

	return tr, nil
}
