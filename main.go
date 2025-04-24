package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/auth"
	"github.com/maty546/secure_payment_service_challenge/controller"
	"github.com/maty546/secure_payment_service_challenge/db"
	"github.com/maty546/secure_payment_service_challenge/repository"
	"github.com/maty546/secure_payment_service_challenge/routes"
	"github.com/maty546/secure_payment_service_challenge/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	db := db.ConnectDB()
	newAccountRepo := repository.NewAccountRepository(db)
	newTransferRepo := repository.NewTransferRepository(db)
	newService := service.NewService(newAccountRepo, newTransferRepo)
	newController := controller.NewController(newService)

	//auth user and pass could be a config, ideally a secret
	newLoginHandler := auth.NewLogin("api", "123")
	routes.RegisterRoutes(r, newController, newLoginHandler)

	r.Run(":8080")
}
