package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/auth"
	"github.com/maty546/secure_payment_service_challenge/controller"
)

func RegisterRoutes(r *gin.Engine, controller controller.ISecurePaymentsController, authHandler auth.IAuthHandler) {

	r.POST("/login", authHandler.GetToken)

	securePaymentsGroup := r.Group("/secure-payments")

	securePaymentsAuthGroup := securePaymentsGroup.Group("/authorized-use")
	securePaymentsAuthGroup.Use(authHandler.GetAuthMiddleware())

	securePaymentsAuthGroup.POST("/transfer", controller.HandleTransferStart)
	securePaymentsAuthGroup.GET("/transfer/:id", controller.HandleTransferGet)
	securePaymentsAuthGroup.GET("/accounts/:id", controller.HandleAccountGet)

	callbackGroup := securePaymentsAuthGroup.Group("/callback")
	callbackGroup.POST("/transfer/result", controller.HandleTransferResultCallback)

	securePaymentsGroupNoAuthGroup := securePaymentsGroup.Group("/no-auth")
	securePaymentsGroupNoAuthGroup.POST("/transfer/timeout/:id", controller.HandleTimeoutCheckForTransfer)
}
