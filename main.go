package main

import (
	"context"
	"log"
	"os"

	"github.com/kinta-mti/mobbe/config"
	"github.com/kinta-mti/mobbe/db"
	"github.com/kinta-mti/mobbe/pushnotification"
	"github.com/kinta-mti/mobbe/service"
	"github.com/kinta-mti/mobbe/ypg"
)

func main() {
	log.Print("[main] v0.0.0:2024-11-06:01")
	ctx := context.Background()
	cfg := config.Load(os.Args[1])
	ypg.Init(cfg.Ypg.ApiKey, cfg.Ypg.SecretKey,
		cfg.Ypg.Apimkey, cfg.Ypg.ApimSecret,
		cfg.Ypg.Path.Uri, cfg.Ypg.Path.AccesToken, cfg.Ypg.Path.Inquiries)
	db.Init(cfg.Database.Name, cfg.Database.User, cfg.Database.Pass)
	pushnotification.Init(ctx)
	service.Init(cfg.Server.Port, ctx)
	//	router := gin.Default()
	//	router.POST("/checkout", endpoint.PostCheckout)
	//	router.POST("/webhook", endpoint.PostWebhook)
	//	router.GET("/helo", endpoint.GetWorld)
	//	router.Run(":" + cfg.Server.Port)
}
