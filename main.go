package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kinta-mti/mobbe/config"
	"github.com/kinta-mti/mobbe/endpoint"
)

var cfg config.Configuration

func main() {
	cfg := config.Load(os.Args[0])
	router := gin.Default()
	router.POST("/checkout", endpoint.PostCheckout)
	router.POST("/webhook", endpoint.PostWebhook)
	router.GET("/helo", endpoint.GetWorld)
	router.Run(":" + cfg.Server.Port)
}
