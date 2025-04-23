package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/service"
)

type ISecurePaymentsController interface {
	GetItem(c *gin.Context)
}

type securePaymentsController struct {
	service service.ISecurePaymentsService
}

func NewController(service service.ISecurePaymentsService) ISecurePaymentsController {
	return securePaymentsController{service}
}

var _ ISecurePaymentsController = (securePaymentsController{})

func (s securePaymentsController) GetItem(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		//fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}
	result, err := s.service.GetItem(c, num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, result)
}
