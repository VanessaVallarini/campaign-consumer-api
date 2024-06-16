package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/api"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/dao"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/listener/handler"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/listener/processor"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/consumer"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/postgres"
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

	// repository
	pool := postgres.CreatePool(ctx, &cfg.Database)
	ownerRepository := dao.NewOwnerRepository(pool)
	slugRepository := dao.NewSlugRepository(pool)
	regionRepository := dao.NewRegionRepository(pool)
	merchantRepository := dao.NewMerchantRepository(pool)
	campaignRepository := dao.NewCampaignRepository(pool)

	// service
	ownerService := service.NewOwnerService(ownerRepository)
	slugService := service.NewSlugService(slugRepository)
	regionService := service.NewRegionService(regionRepository)
	merchantService := service.NewMerchantService(merchantRepository)
	campaignService := service.NewCampaignService(campaignRepository)

	// processor
	slugProcessor := processor.NewSlugProcessor(slugService)
	regionProcessor := processor.NewRegionProcessor(regionService)
	merchantProcessor := processor.NewMerchantProcessor(merchantService)
	campaignProcessor := processor.NewCampaignProcessor(campaignService)

	// handler
	ownerHandler := handler.MakeOwnerEventHandler(ownerService)
	slugHandler := handler.MakeSlugEventHandler(slugProcessor)
	regionHandler := handler.MakeRegionEventHandler(regionProcessor)
	merchantHandler := handler.MakeMerchantEventHandler(merchantProcessor)
	campaignHandler := handler.MakeCampaignEventHandler(campaignProcessor)

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
		easyzap.Info(ctx, "Starting http worker server at "+cfg.ServerHost)
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
		easyzap.Info(ctx, "Starting metadata worker server at "+cfg.MetaHost)
		err := meta.Start(cfg.MetaHost)
		easyzap.Fatal(ctx, err, "failed to start meta server")
	}()

	// listens for system signals to gracefully shutdown
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	switch <-signalChannel {
	case os.Interrupt:
		easyzap.Info(context.Background(), "Received SIGINT, stopping...")
	case syscall.SIGTERM:
		easyzap.Info(context.Background(), "Received SIGTERM, stopping...")
	}

}
