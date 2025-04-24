package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/auth"
	"github.com/maty546/secure_payment_service_challenge/controller"
)

func RegisterRoutes(r *gin.Engine, controller controller.ISecurePaymentsController, authHandler auth.IAuthHandler) {

	r.POST("/login", authHandler.GetToken)

	securePaymentsGroup := r.Group("/secure-payments")
	securePaymentsGroup.Use(authHandler.GetAuthMiddleware())

	securePaymentsGroup.POST("/transfer", controller.HandleTransferStart)
	securePaymentsGroup.GET("/transfer/:id", controller.HandleTransferGet)
	securePaymentsGroup.GET("/accounts/:id", controller.HandleAccountGet)

	callbackGroup := securePaymentsGroup.Group("/callback")
	callbackGroup.POST("/transfer/result", controller.HandleTransferResultCallback)
}
