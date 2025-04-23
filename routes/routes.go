package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/controller"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", controller.Ping)
}
