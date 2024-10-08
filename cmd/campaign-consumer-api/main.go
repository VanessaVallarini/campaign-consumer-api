package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/api"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/dao"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/listener/handler"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/consumer"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/postgres"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	easyzap "github.com/lockp111/go-easyzap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoTracer "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(echoTracer.Middleware())
	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST},
	}))

	timeLocation, err := time.LoadLocation(cfg.TimeLocation)
	if err != nil {
		easyzap.Fatal(ctx, err, "failed to load timeLocation")
	}

	// dao
	pool := postgres.CreatePool(ctx, &cfg.Database)
	ownerDao := dao.NewOwnerDao(pool)
	slugDao := dao.NewSlugDao(pool)
	regionDao := dao.NewRegionDao(pool)
	merchantDao := dao.NewMerchantDao(pool)
	campaignDao := dao.NewCampaignDao(pool)
	campaignHistoryDao := dao.NewCampaignHistoryDao(pool)
	slugHistoryDao := dao.NewSlugHistoryDao(pool)
	regionHistoryDao := dao.NewRegionHistoryDao(pool)
	spentDao := dao.NewSpentDao(pool)
	transactionManager := transaction.NewTransactionManager(pool)

	// service
	ownerService := service.NewOwnerService(ownerDao)
	slugService := service.NewSlugService(slugDao, slugHistoryDao, transactionManager)
	regionService := service.NewRegionService(regionDao, regionHistoryDao, transactionManager)
	merchantService := service.NewMerchantService(merchantDao)
	spentService := service.NewSpentService(spentDao)
	bucketService := service.NewBucketService(timeLocation)
	campaignService := service.NewCampaignService(campaignDao, campaignHistoryDao, spentService, bucketService, transactionManager)

	// handler
	ownerHandler := handler.MakeOwnerEventHandler(ownerService)
	slugHandler := handler.MakeSlugEventHandler(slugService)
	regionHandler := handler.MakeRegionEventHandler(regionService)
	merchantHandler := handler.MakeMerchantEventHandler(merchantService)
	campaignHandler := handler.MakeCampaignEventHandler(campaignService)

	// client
	ownerSrClient := client.NewSchemaRegistry(cfg.KafkaOwner)
	slugSrClient := client.NewSchemaRegistry(cfg.KafkaSlug)
	regionSrClient := client.NewSchemaRegistry(cfg.KafkaRegion)
	merchantSrClient := client.NewSchemaRegistry(cfg.KafkaMerchant)
	campaignSrClient := client.NewSchemaRegistry(cfg.KafkaCampaign)

	//consumer
	ownerConsumer := consumer.NewConsumer(ctx, cfg.KafkaOwner, ownerSrClient, ownerHandler)
	go ownerConsumer.ConsumerStart(cfg.KafkaOwner)

	slugConsumer := consumer.NewConsumer(ctx, cfg.KafkaSlug, slugSrClient, slugHandler)
	go slugConsumer.ConsumerStart(cfg.KafkaSlug)

	regionConsumer := consumer.NewConsumer(ctx, cfg.KafkaRegion, regionSrClient, regionHandler)
	go regionConsumer.ConsumerStart(cfg.KafkaRegion)

	merchantConsumer := consumer.NewConsumer(ctx, cfg.KafkaMerchant, merchantSrClient, merchantHandler)
	go merchantConsumer.ConsumerStart(cfg.KafkaMerchant)

	campaignConsumer := consumer.NewConsumer(ctx, cfg.KafkaCampaign, campaignSrClient, campaignHandler)
	go campaignConsumer.ConsumerStart(cfg.KafkaCampaign)

	// Start HTTP server
	go func() {
		easyzap.Info(ctx, "starting http worker server at "+cfg.ServerHost)
		err := server.Start(cfg.ServerHost)
		easyzap.Fatal(ctx, err, "failed to start server")
	}()

	meta := echo.New()
	meta.HideBanner = true
	meta.HidePort = true

	meta.GET("/prometheus", echo.WrapHandler(promhttp.Handler()))

	api.NewHealthCheck().Register(meta)

	// starts meta application
	go func() {
		easyzap.Info(ctx, "starting metadata worker server at "+cfg.MetaHost)
		err := meta.Start(cfg.MetaHost)
		easyzap.Fatal(ctx, err, "failed to start meta server")
	}()

	// listens for system signals to gracefully shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	switch <-signalChannel {
	case os.Interrupt:
		easyzap.Info(context.Background(), "received SIGINT, stopping...")
	case syscall.SIGTERM:
		easyzap.Info(context.Background(), "received SIGTERM, stopping...")
	}

}
