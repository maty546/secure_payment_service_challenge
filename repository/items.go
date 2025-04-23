package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"gorm.io/gorm"
)

type ISecurePaymentsRepository interface {
	GetItem(c *gin.Context, id int) (models.Item, error)
}

type securePaymentsRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) securePaymentsRepository {
	return securePaymentsRepository{db}
}

var _ ISecurePaymentsRepository = (securePaymentsRepository{})

func (r securePaymentsRepository) GetItem(c *gin.Context, id int) (models.Item, error) {
	var item models.Item
	if result := r.db.First(&item, id); result.Error != nil {
		log.Error(fmt.Sprintf("Repository - %s", result.Error.Error()))
		return models.Item{}, result.Error
	}

	return item, nil
}
