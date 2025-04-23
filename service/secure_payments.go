package service

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"github.com/maty546/secure_payment_service_challenge/repository"
)

type ISecurePaymentsService interface {
	StartTransfer(c *gin.Context, transfer models.Transfer) (models.Transfer, error)
	GetTransferByID(c *gin.Context, transferID uint) (models.Transfer, error)
	GetAccountByID(c *gin.Context, accountID string) (models.Account, error)
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
	errTransferToSameAcc                = "cant transfer to the same account"
	errExternalToExternal               = "cant transfer if both accounts are external"
	errInsuficientFunds                 = "the giving account doesnt have enough funds for this transfer"
)

//todo transform error vars into actual errors

func isExternalAccount(accountID string) bool {
	return strings.Contains(accountID, "ext")
}

// checks if account exists and if it has sufficient funds
func (s securePaymentsService) accountAbleToTransfer(c *gin.Context, accountID string, transfer models.Transfer) bool {
	//get acc balance
	acc, errGettingAcc := s.accountsRepo.GetByID(c, accountID)
	if errGettingAcc != nil {
		log.Error(fmt.Sprintf("securePaymentsService | accountAbleToTransfer err - %s", errGettingAcc.Error()))
		return false
	}

	pendingPaymentsAmount, errGettingPendingPaymentsAmount := s.transfersRepo.GetPendingPaymentsAmountForAccount(c, accountID)
	if errGettingPendingPaymentsAmount != nil {
		log.Error(fmt.Sprintf("securePaymentsService | accountAbleToTransfer err - %s", errGettingPendingPaymentsAmount.Error()))
		return false
	}
	//todo remove
	log.Info(fmt.Sprintf("balance %d, pending %d, transfer %d", acc.Balance, pendingPaymentsAmount, transfer.Amount))
	return acc.Balance-pendingPaymentsAmount >= transfer.Amount
}

func (s securePaymentsService) StartTransfer(c *gin.Context, transfer models.Transfer) (models.Transfer, error) {

	if transfer.Amount == 0 {
		return models.Transfer{}, nil
	}

	//if its a same acc transfer or external to external, do nothing
	if transfer.ToAccountID == transfer.FromAccountID {
		log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err - %s", errTransferToSameAcc))
		return models.Transfer{}, errors.New(errTransferToSameAcc)
	}

	fromIsExternal := isExternalAccount(transfer.FromAccountID)
	toIsExternal := isExternalAccount(transfer.ToAccountID)

	if fromIsExternal && toIsExternal {
		log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err - %s", errExternalToExternal))
		return models.Transfer{}, errors.New(errExternalToExternal)
	}
	//validation for case internal->external
	if !fromIsExternal && !s.accountAbleToTransfer(c, transfer.FromAccountID, transfer) {
		log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err - %s", errInsuficientFunds))
		return models.Transfer{}, errors.New(errInsuficientFunds)
	}

	if !toIsExternal {
		// check that to account exists
		_, errGettingToAcc := s.accountsRepo.GetByID(c, transfer.ToAccountID)
		if errGettingToAcc != nil {
			log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err - %s", errGettingToAcc.Error()))
			return models.Transfer{}, fmt.Errorf(errCouldNotObtainDestinationAccount, errGettingToAcc)
		}
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

func (s securePaymentsService) GetAccountByID(c *gin.Context, accountID string) (models.Account, error) {
	acc, err := s.accountsRepo.GetByID(c, accountID)
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsService | GetAccountByID err - %s", err.Error()))
		return models.Account{}, fmt.Errorf(errCouldNotGetAcc, err)
	}

	return acc, nil
}
