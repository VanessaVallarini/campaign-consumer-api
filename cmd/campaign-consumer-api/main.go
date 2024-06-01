package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/api"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/consumer"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/handler"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/processor"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/repository"
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
	pool := repository.CreatePool(ctx, &cfg.Database)
	ownerRepository := repository.NewOwnerRepository(pool)
	slugRepository := repository.NewSlugRepository(pool)
	regionRepository := repository.NewRegionRepository(pool)

	// service
	ownerService := service.NewOwnerService(ownerRepository)
	slugService := service.NewSlugService(slugRepository)
	regionService := service.NewRegionService(regionRepository)

	// processor
	ownerProcessor := processor.NewOwnerProcessor(ownerService)
	slugProcessor := processor.NewSlugProcessor(slugService)
	regionProcessor := processor.NewRegionProcessor(regionService)

	// handler
	ownerHandler := handler.MakeOwnerEventHandler(ownerProcessor)
	slugHandler := handler.MakeSlugEventHandler(slugProcessor)
	regionHandler := handler.MakeRegionEventHandler(regionProcessor)

	// client
	ownerSrClient := client.NewSchemaRegistry(cfg.KafkaOwner)
	slugSrClient := client.NewSchemaRegistry(cfg.KafkaSlug)
	regionSrClient := client.NewSchemaRegistry(cfg.KafkaRegion)

	//consumer
	ownerConsumer := consumer.NewConsumer(ctx, cfg.KafkaOwner, ownerSrClient, ownerHandler)
	go ownerConsumer.ConsumerStart(cfg.KafkaOwner)

	slugConsumer := consumer.NewConsumer(ctx, cfg.KafkaSlug, slugSrClient, slugHandler)
	go slugConsumer.ConsumerStart(cfg.KafkaSlug)

	regionConsumer := consumer.NewConsumer(ctx, cfg.KafkaRegion, regionSrClient, regionHandler)
	go regionConsumer.ConsumerStart(cfg.KafkaRegion)

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
