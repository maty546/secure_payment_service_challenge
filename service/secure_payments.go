package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"github.com/maty546/secure_payment_service_challenge/repository"
)

type ISecurePaymentsService interface {
	StartTransfer(c *gin.Context, transfer models.Transfer) (models.Transfer, error)
}

type securePaymentsService struct {
	accountsRepo  repository.IAccountRepository
	transfersRepo repository.ITransferRepository
}

func NewService(accountsRepo repository.IAccountRepository, transfersRepo repository.ITransferRepository) securePaymentsService {
	return securePaymentsService{accountsRepo, transfersRepo}
}

var _ ISecurePaymentsService = (securePaymentsService{})

var errCouldNotObtainDestinationAccount = "could not obtain destination account - %#v"
var errCouldNotObtainOriginAccount = "could not obtain origin account - %#v"
var errCouldNotStartTransfer = "could not start transfer - %#v"

func (s securePaymentsService) StartTransfer(c *gin.Context, transfer models.Transfer) (models.Transfer, error) {
	//check that to and from exist
	_, errGettingToAcc := s.accountsRepo.GetByID(c, transfer.ToAccountID)
	if errGettingToAcc != nil {
		log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err - %s", errGettingToAcc.Error()))
		return models.Transfer{}, fmt.Errorf(errCouldNotObtainDestinationAccount, errGettingToAcc)
	}

	_, errGettingFromAcc := s.accountsRepo.GetByID(c, transfer.FromAccountID)
	if errGettingFromAcc != nil {
		log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err - %s", errGettingFromAcc.Error()))
		return models.Transfer{}, fmt.Errorf(errCouldNotObtainOriginAccount, errGettingFromAcc)
	}
	//save transfer
	transfer.Status = models.TRANSFER_STATUS_PENDING
	savedTransfer, errSaving := s.transfersRepo.Save(c, transfer)
	if errSaving != nil {
		log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err - %s", errSaving.Error()))
		return models.Transfer{}, fmt.Errorf(errCouldNotStartTransfer, errSaving)
	}

	return savedTransfer, nil
}
