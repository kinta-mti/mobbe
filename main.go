package main

import (
	"log"
	"os"

	//	"github.com/gin-gonic/gin"
	"github.com/kinta-mti/mobbe/config"
	"github.com/kinta-mti/mobbe/db"
	"github.com/kinta-mti/mobbe/ypg"
	//	"github.com/kinta-mti/mobbe/endpoint"
)

func main() {
	cfg := config.Load(os.Args[1])
	ypg.Init(cfg.Ypg.ApiKey, cfg.Ypg.SecretKey,
		cfg.Ypg.Apimkey, cfg.Ypg.ApimSecret,
		cfg.Ypg.Path.Uri, cfg.Ypg.Path.AccesToken, cfg.Ypg.Path.Inquiries)
	db.Init(cfg.Database.Name, cfg.Database.User, cfg.Database.Pass)
	log.Print("[main]" + cfg.Server.Port)
	//	endpoint.Init(cfg)
	//	router := gin.Default()
	//	router.POST("/checkout", endpoint.PostCheckout)
	//	router.POST("/webhook", endpoint.PostWebhook)
	//	router.GET("/helo", endpoint.GetWorld)
	//	router.Run(":" + cfg.Server.Port)
}
