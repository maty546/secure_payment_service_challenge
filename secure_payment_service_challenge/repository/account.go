package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	GetByID(c *gin.Context, id string) (models.Account, error)
	Save(c *gin.Context, acc models.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) accountRepository {
	return accountRepository{db}
}

var _ IAccountRepository = (accountRepository{})

func (r accountRepository) GetByID(c *gin.Context, id string) (models.Account, error) {
	var item models.Account
	if result := r.db.First(&item, "id = ?", id); result.Error != nil {
		log.Error(fmt.Sprintf("accountRepository | GetByID err - %s", result.Error.Error()))
		return models.Account{}, result.Error
	}

	return item, nil
}

func (r accountRepository) Save(c *gin.Context, acc models.Account) error {
	if result := r.db.Save(&acc); result.Error != nil {
		log.Error(fmt.Sprintf("accountRepository | Save err - %s", result.Error.Error()))
		//todo type this error
		return result.Error
	}

	return nil
}
