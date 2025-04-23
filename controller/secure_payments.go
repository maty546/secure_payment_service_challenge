package controller

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/service"
)

type ISecurePaymentsController interface {
	HandleTransferStart(c *gin.Context)
	HandleTransferGet(c *gin.Context)
	HandleAccountGet(c *gin.Context)
}

type securePaymentsController struct {
	service service.ISecurePaymentsService
}

func NewController(service service.ISecurePaymentsService) ISecurePaymentsController {
	return securePaymentsController{service}
}

var _ ISecurePaymentsController = (securePaymentsController{})

func (s securePaymentsController) HandleTransferStart(c *gin.Context) {

	var requestBody HandleTransferStartRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Error(fmt.Sprintf("securePaymentsController | HandleTransferStart err - %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transfer := requestBody.parseIntoTransferModel()

	result, err := s.service.StartTransfer(c, transfer)
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsController | HandleTransferStart err - %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s securePaymentsController) HandleTransferGet(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		log.Error(fmt.Sprintf("securePaymentsController | HandleTransferGet err - %s", err.Error()))
		return
	}

	transfer, err := s.service.GetTransferByID(c, uint(id))
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsController | HandleTransferGet err - %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (s securePaymentsController) HandleAccountGet(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		log.Error(fmt.Sprintf("securePaymentsController | HandleAccountGet err - %s", err.Error()))
		return
	}

	account, err := s.service.GetAccountByID(c, uint(id))
	if err != nil {
		log.Error(fmt.Sprintf("securePaymentsController | HandleAccountGet err - %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}
