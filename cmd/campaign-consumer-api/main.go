package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/api"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/dao"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/listener/handler"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/listener/processor"
	kafkaClient "github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/consumer"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/postgres"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/redis"
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

	// client
	redisClient := redis.NewRedisDataClient(cfg.RedisConfig.RWAddress, cfg.RedisConfig.ROAddress)
	redisClient.Ping(ctx)
	slugRedisValidator := service.NewRedisValidator(redisClient, cfg.RedisConfig.SlugTTL)
	regionRedisValidator := service.NewRedisValidator(redisClient, cfg.RedisConfig.RegionTTL)
	addressApiClient := client.NewAddressClient(cfg.AddressApi)

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
	ledgerDao := dao.NewLedgerDao(pool)
	transactionManager := transaction.NewTransactionManager(pool)

	// service
	ownerService := service.NewOwnerService(ownerDao)
	slugService := service.NewSlugService(slugDao, slugHistoryDao, transactionManager, slugRedisValidator)
	regionService := service.NewRegionService(regionDao, regionHistoryDao, transactionManager, regionRedisValidator)
	merchantService := service.NewMerchantService(merchantDao)
	bucketService := service.NewBucketService(timeLocation)
	ledgerServce := service.NewLedgerService(ledgerDao)
	spentService := service.NewSpentService(spentDao, transactionManager, ledgerServce)
	campaignService := service.NewCampaignService(campaignDao, campaignHistoryDao, spentService, bucketService, transactionManager)

	// processor
	spentProcessor := processor.NewSpentProcessor(spentService, campaignService, merchantService, slugService, regionService, ledgerServce, bucketService, addressApiClient)

	// handler
	ownerHandler := handler.MakeOwnerEventHandler(ownerService)
	slugHandler := handler.MakeSlugEventHandler(slugService)
	regionHandler := handler.MakeRegionEventHandler(regionService)
	merchantHandler := handler.MakeMerchantEventHandler(merchantService)
	campaignHandler := handler.MakeCampaignEventHandler(campaignService)
	spentHandler := handler.MakeSpentEventHandler(spentProcessor)

	ownerSrClient := kafkaClient.NewSchemaRegistry(cfg.KafkaOwner)
	slugSrClient := kafkaClient.NewSchemaRegistry(cfg.KafkaSlug)
	regionSrClient := kafkaClient.NewSchemaRegistry(cfg.KafkaRegion)
	merchantSrClient := kafkaClient.NewSchemaRegistry(cfg.KafkaMerchant)
	campaignSrClient := kafkaClient.NewSchemaRegistry(cfg.KafkaCampaign)
	spentSrClient := kafkaClient.NewSchemaRegistry(cfg.KafkaSpent)

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

	spentConsumer := consumer.NewConsumer(ctx, cfg.KafkaSpent, spentSrClient, spentHandler)
	go spentConsumer.ConsumerStart(cfg.KafkaSpent)

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
