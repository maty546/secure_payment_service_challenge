package service

import (
	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"github.com/maty546/secure_payment_service_challenge/repository"
)

type ISecurePaymentsService interface {
	GetItem(c *gin.Context, id int) (models.Item, error)
}

type securePaymentsService struct {
	repo repository.ISecurePaymentsRepository
}

func NewService(repo repository.ISecurePaymentsRepository) securePaymentsService {
	return securePaymentsService{repo}
}

var _ ISecurePaymentsService = (securePaymentsService{})

func (s securePaymentsService) GetItem(c *gin.Context, id int) (models.Item, error) {
	return s.repo.GetItem(c, id)
}
