package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maty546/secure_payment_service_challenge/routes"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)

	r.Run(":8080") // start server on port 8080
}
