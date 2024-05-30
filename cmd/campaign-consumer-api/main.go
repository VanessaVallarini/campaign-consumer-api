package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/api"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
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

	// clients
	kafkaOwnerClient := client.NewKafkaClient(ctx, cfg.KafkaOwner)
	schemaRegistryOwnerClient := client.NewSchemaRegistry(cfg.KafkaOwner)

	kafkaSlugClient := client.NewKafkaClient(ctx, cfg.KafkaSlug)
	schemaRegistrySlugClient := client.NewSchemaRegistry(cfg.KafkaSlug)

	kafkaMerchantClient := client.NewKafkaClient(ctx, cfg.KafkaMerchant)
	schemaRegistryMerchantClient := client.NewSchemaRegistry(cfg.KafkaMerchant)

	kafkaCampaignClient := client.NewKafkaClient(ctx, cfg.KafkaCampaign)
	schemaRegistryCampaignClient := client.NewSchemaRegistry(cfg.KafkaCampaign)

	KafkaClickImpressionClient := client.NewKafkaClient(ctx, cfg.KafkaClickImpression)
	schemaRegistryClickImpressionClient := client.NewSchemaRegistry(cfg.KafkaClickImpression)

	fmt.Println(kafkaOwnerClient)
	fmt.Println(schemaRegistryOwnerClient)
	fmt.Println(kafkaSlugClient)
	fmt.Println(schemaRegistrySlugClient)
	fmt.Println(kafkaMerchantClient)
	fmt.Println(schemaRegistryMerchantClient)
	fmt.Println(kafkaCampaignClient)
	fmt.Println(schemaRegistryCampaignClient)
	fmt.Println(KafkaClickImpressionClient)
	fmt.Println(schemaRegistryClickImpressionClient)

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
