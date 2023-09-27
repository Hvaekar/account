package commands

import (
	"fmt"
	"github.com/Hvaekar/med-account/cmd/account/handler"
	"github.com/Hvaekar/med-account/config"
	"github.com/Hvaekar/med-account/internal/account"
	"github.com/Hvaekar/med-account/pkg/amazon"
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/broker/kafka"
	"github.com/Hvaekar/med-account/pkg/logger"
	"github.com/Hvaekar/med-account/pkg/server"
	"github.com/Hvaekar/med-account/pkg/server/ginmiddleware"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

func Serve() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Serve HTTP requests",
		Action: func(ctx *cli.Context) error {
			// get configs
			cfg, err := config.Get("config")
			if err != nil {
				return fmt.Errorf("get config: %w", err)
			}

			// create logger
			log := logger.Get(cfg.Logger.LoggerName)
			if err := log.Init(cfg); err != nil {
				return fmt.Errorf("init logger: %w", err)
			}

			// aws session
			aws := amazon.NewAWS(&cfg.AWS)
			awsSession := aws.CreateSession()

			// create app
			router, db, _, _, err := CreateApp(
				cfg,
				log,
				aws,
				s3manager.NewUploader(awsSession),
				s3manager.NewDownloader(awsSession),
			)
			if err != nil {
				return fmt.Errorf("create app: %w", err)
			}
			defer db.Close()

			if ctx.String("space") != "docker" {
				if err := db.Migrate(); err != nil {
					return err
				}
			}

			errCh := make(chan error)

			// run HTTP server
			svr := server.NewServer(router, cfg, log)

			go func() {
				if err := svr.Run(); err != nil {
					errCh <- fmt.Errorf("run server: %w", err)
				}
			}()
			defer svr.Shutdown()

			//// consumer message broker
			//go func() {
			//	if err := mb.Consume(); err != nil {
			//		errCh <- fmt.Errorf("consumer message broker: %w", err)
			//	}
			//}()

			select {
			case <-ctx.Done():
				return nil
			case <-errCh:
				return <-errCh
			}
		},
	}
}

func CreateApp(
	cfg *config.Config,
	log logger.Logger,
	aws *amazon.AWS,
	s3Uploader s3manageriface.UploaderAPI,
	s3Downloader s3manageriface.DownloaderAPI,
) (*gin.Engine, storage.Storage, *amazon.S3, broker.MessageBroker, error) {
	// connect to postgres storage
	psqlStorage := postgres.NewPostgres(&cfg.Postgres)
	if err := psqlStorage.Connect(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("connect postgres: %w", err)
	}

	// create message broker
	kafkaMB := kafka.NewMessageBroker(&cfg.Kafka, log)
	//if err := broker.AddConsumer(); err != nil {
	//	return nil, fmt.Errorf("kafka add consumer: %w", err)
	//}
	if err := kafkaMB.AddProducer(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("kafka add producer: %w", err)
	}

	// aws s3
	awsS3 := amazon.NewS3(aws, s3Uploader, s3Downloader)

	// create router and handlers
	router := gin.New()

	psql := account.NewPostgresStorage(psqlStorage)

	basicH := handler.NewBasicHandler(log, cfg, psql, awsS3, kafkaMB)
	authH := handler.NewAuthHandler(basicH)
	accountH := handler.NewAccountHandler(basicH)
	fileH := handler.NewFileHandler(basicH)
	emailH := handler.NewEmailHandler(basicH)
	phoneH := handler.NewPhoneHandler(basicH)
	addressH := handler.NewAddressHandler(basicH)
	langH := handler.NewLanguageHandler(basicH)
	profileH := handler.NewProfileHandler(basicH)
	patientH := handler.NewPatientHandler(basicH)
	metalComponentH := handler.NewMetalComponentHandler(basicH)
	adminH := handler.NewPatientAdminHandler(basicH)
	specialistH := handler.NewSpecialistHandler(basicH)
	specializationH := handler.NewSpecializationHandler(basicH)
	educationH := handler.NewEducationHandler(basicH)
	experienceH := handler.NewExperienceHandler(basicH)
	associationH := handler.NewAssociationHandler(basicH)
	patentH := handler.NewPatentHandler(basicH)
	publicationH := handler.NewPublicationHandler(basicH)

	router.Use(
		gin.Recovery(),
		ginmiddleware.Cors(),
		ginmiddleware.Logger(log, "/health", "/metrics"),
		ginmiddleware.PrometheusMetrics("/health", "/metrics"),
		//ginmiddleware.AWSSession(aws.GetSession()),
	)

	router.GET("/health", server.Health)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	authH.InitRoutes(router)
	arg := accountH.InitRoutes(router)
	fileH.InitRoutes(arg)
	emailH.InitRoutes(arg)
	phoneH.InitRoutes(arg)
	addressH.InitRoutes(arg)
	langH.InitRoutes(arg)
	profileH.InitRoutes(arg)
	prg := patientH.InitRoutes(router)
	metalComponentH.InitRoutes(prg)
	adminH.InitRoutes(prg)
	srg := specialistH.InitRoutes(router)
	specializationH.InitRoutes(srg)
	educationH.InitRoutes(srg)
	experienceH.InitRoutes(srg)
	associationH.InitRoutes(srg)
	patentH.InitRoutes(srg)
	publicationH.InitRoutes(srg)

	pprof.Register(router)

	return router, psqlStorage, awsS3, kafkaMB, nil
}
