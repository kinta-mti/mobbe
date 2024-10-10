package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kinta-mti/mobbe/config"
)

var cfg config.Configuration

func main() {
	cfg := config.Load(os.Args[0])
	router := gin.Default()
	router.POST("/checkout", endpoint.postCheckout)
	router.POST("/webhook", endpoint.postWebhook)
	router.GET("/helo", endpoint.getWorld)
	router.Run(":" + cfg.Server.Port)
}
