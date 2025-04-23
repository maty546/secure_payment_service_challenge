package controller

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/service"
)

type ISecurePaymentsController interface {
	HandleTransferStart(c *gin.Context)
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
