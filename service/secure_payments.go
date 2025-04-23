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
	GetTransferByID(c *gin.Context, transferID uint) (models.Transfer, error)
	GetAccountByID(c *gin.Context, accountID uint) (models.Account, error)
}

type securePaymentsService struct {
	accountsRepo  repository.IAccountRepository
	transfersRepo repository.ITransferRepository
}

func NewService(accountsRepo repository.IAccountRepository, transfersRepo repository.ITransferRepository) securePaymentsService {
	return securePaymentsService{accountsRepo, transfersRepo}
}

var _ ISecurePaymentsService = (securePaymentsService{})

var (
	errCouldNotObtainDestinationAccount = "could not obtain destination account - %#v"
	errCouldNotObtainOriginAccount      = "could not obtain origin account - %#v"
	errCouldNotStartTransfer            = "could not start transfer - %#v"
	errCouldNotGetAcc                   = "could not get account - %#v"
	errCouldNotGetTr                    = "could not get transfer - %#v"
)

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

func (s securePaymentsService) GetTransferByID(c *gin.Context, transferID uint) (models.Transfer, error) {
	transfer, err := s.transfersRepo.GetByID(c, transferID)
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsService | GetTransferByID err - %s", err.Error()))
		return models.Transfer{}, fmt.Errorf(errCouldNotGetTr, err)
	}

	return transfer, nil
}

func (s securePaymentsService) GetAccountByID(c *gin.Context, accountID uint) (models.Account, error) {
	acc, err := s.accountsRepo.GetByID(c, accountID)
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsService | GetAccountByID err - %s", err.Error()))
		return models.Account{}, fmt.Errorf(errCouldNotGetAcc, err)
	}

	return acc, nil
}
