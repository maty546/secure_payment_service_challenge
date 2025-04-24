package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/models"
	"github.com/maty546/secure_payment_service_challenge/repository"
)

type ISecurePaymentsService interface {
	StartTransfer(c *gin.Context, transfer models.Transfer) (models.Transfer, error)
	GetTransferByID(c *gin.Context, transferID uint) (models.Transfer, error)
	GetAccountByID(c *gin.Context, accountID string) (models.Account, error)
	HandleTransferResultCallback(c *gin.Context, transferID uint, status models.TransferStatus) error
	HandleTimeoutCheckForTransfer(c *gin.Context, transferID uint) error
}

type securePaymentsService struct {
	accountsRepo                repository.IAccountRepository
	transfersRepo               repository.ITransferRepository
	asynqClient                 *asynq.Client
	timeoutCheckForTransferAddr string
}

func NewService(accountsRepo repository.IAccountRepository, transfersRepo repository.ITransferRepository, asynqClientAddr string, timeoutCheckAddr string) securePaymentsService {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: asynqClientAddr})

	return securePaymentsService{accountsRepo, transfersRepo, client, timeoutCheckAddr}
}

var _ ISecurePaymentsService = (securePaymentsService{})

var (
	errCouldNotObtainDestinationAccount = "could not obtain destination account - %#v"
	errCouldNotObtainOriginAccount      = "could not obtain origin account - %#v"
	errCouldNotPerformTransfer          = "could not perform transfer - %#v"
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
		return models.Transfer{}, fmt.Errorf(errCouldNotPerformTransfer, errSaving)
	}

	//add async task to workqueue for checking timeouts
	url := s.timeoutCheckForTransferAddr + strconv.FormatUint(uint64(savedTransfer.ID), 10)
	//url := s.timeoutCheckForTransferAddr
	task := asynq.NewTask("http:call", []byte(fmt.Sprintf(`{"url":"%s"}`, url)))

	_, err := s.asynqClient.Enqueue(task, asynq.ProcessIn(30*time.Second))
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsService | StartTransfer err enqueuing timeout check - %s", err.Error()))
		//should add recovery mechanic
		return models.Transfer{}, err
	}
	log.Info(fmt.Sprintf("securePaymentsService | StartTransfer just enqueued task successfully with url %s", url))
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

func (s securePaymentsService) HandleTransferResultCallback(c *gin.Context, transferID uint, status models.TransferStatus) error {

	if status != models.TRANSFER_STATUS_COMPLETED && status != models.TRANSFER_STATUS_FAILED {
		log.Error("securePaymentsService | HandleTransferResultCallback err - invalid status received")
		return errors.New("invalid status received")
	}

	transfer, err := s.transfersRepo.GetByID(c, transferID)
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsService | HandleTransferResultCallback err - %s", err.Error()))
		return errors.New("could not obtain transfer")
	}

	if transfer.Status != models.TRANSFER_STATUS_PENDING {
		log.Error("securePaymentsService | HandleTransferResultCallback err - transaction was already processed")
		return errors.New("transaction was already processed")
	}

	//CASE FAILED
	if status == models.TRANSFER_STATUS_FAILED {
		s.transfersRepo.SetStatus(c, transferID, models.TRANSFER_STATUS_FAILED)
		return nil
	}

	//CASE COMPLETED

	//todo this behaviour should be in the consumed object
	fromIsExternal := isExternalAccount(transfer.FromAccountID)
	toIsExternal := isExternalAccount(transfer.ToAccountID)

	//case internal to internal
	if !fromIsExternal && !toIsExternal {
		err := s.transfersRepo.CompleteInternalTransfer(c, transfer)
		if err != nil {
			log.Error(fmt.Sprintf("securePaymentsService | HandleTransferResultCallback err - %s", err.Error()))
			return fmt.Errorf(errCouldNotPerformTransfer, err)
		}
		return nil
	}

	//case external to internal
	if fromIsExternal {
		err := s.transfersRepo.ReceiveExternalPayment(c, transfer.ToAccountID, transfer)
		if err != nil {
			log.Error(fmt.Sprintf("securePaymentsService | HandleTransferResultCallback err - %s", err.Error()))
			return fmt.Errorf(errCouldNotPerformTransfer, err)
		}
		return nil
	}

	//case internal to external
	err = s.transfersRepo.MakeExternalPayment(c, transfer.FromAccountID, transfer)
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsService | HandleTransferResultCallback err - %s", err.Error()))
		return fmt.Errorf(errCouldNotPerformTransfer, err)
	}
	return nil

}

func (s securePaymentsService) HandleTimeoutCheckForTransfer(c *gin.Context, transferID uint) error {
	transfer, err := s.transfersRepo.GetByID(c, transferID)
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsService | HandleTimeoutCheckForTransfer err - %s", err.Error()))
		return fmt.Errorf(errCouldNotGetTr, err)
	}

	if transfer.Status == models.TRANSFER_STATUS_PENDING {
		err = s.transfersRepo.SetStatus(c, transferID, models.TRANSFER_STATUS_TIMEOUT)
		if err != nil {
			log.Error(fmt.Sprintf("securePaymentsService | HandleTimeoutCheckForTransfer err - %s", err.Error()))
		}
	}

	return err
}
