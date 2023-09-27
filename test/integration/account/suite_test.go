package account

import (
	"context"
	"github.com/Hvaekar/med-account/cmd/account/commands"
	"github.com/Hvaekar/med-account/config"
	"github.com/Hvaekar/med-account/pkg/amazon"
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/client"
	"github.com/Hvaekar/med-account/pkg/dockertest"
	"github.com/Hvaekar/med-account/pkg/dockertest/kafkatest"
	"github.com/Hvaekar/med-account/pkg/dockertest/postgrestest"
	"github.com/Hvaekar/med-account/pkg/dockertest/zookeepertest"
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/logger"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
)

var truncateTables = []string{
	"accounts",
	"account_files",
	"account_emails",
	"account_phones",
	"account_addresses",
	"account_languages",
	"patient_profiles",
	"accounts_patient_profiles",
	"patient_disability_files",
	"patient_metal_components",
	"specialist_profiles",
	"specialist_specializations",
	"specialist_cures_diseases",
	"specialist_services",
	"specialist_educations",
	"specialist_education_files",
	"specialist_experiences",
	"specialist_experience_specializations",
	"specialist_associations",
	"specialist_patents",
	"specialist_publication_links",
}

type TestSuite struct {
	suite.Suite
	ctx    context.Context
	cm     *dockertest.ContainerManager
	cfg    *config.Config
	db     storage.Storage
	s3     *amazon.S3
	mb     broker.MessageBroker
	token  *model.Token
	client *client.HTTPClient

	uploaderAPIMock   *amazon.MockUploaderAPI
	downloaderAPIMock *amazon.MockDownloaderAPI
}

func (s *TestSuite) SetupSuite() {
	err := godotenv.Load("../../../.env")
	s.Require().NoError(err)

	// get configs
	cfg, err := config.Get("config", "../../../config")
	s.Require().NoError(err)

	// change log level to debug
	cfg.Logger.Level = "debug"

	s.cfg = cfg

	// init logger
	log := logger.Get("zap")
	err = log.Init(cfg)
	s.Require().NoError(err)

	// add context
	s.ctx = context.Background()

	// create container manager
	cm := dockertest.NewContainerManager(log)
	err = cm.CreatePool()
	s.Require().NoError(err)
	err = cm.AddNetwork("test-account")
	s.Require().NoError(err)

	// create objects for containers
	psqlContainer := postgrestest.NewDefaultContainer(cm.GetNetwork())
	zooContainer := zookeepertest.NewDefaultContainer(cm.GetNetwork())
	kafkaContainer := kafkatest.NewDefaultContainer(cm.GetNetwork())

	cm.AddContainer(psqlContainer)
	cm.AddContainer(zooContainer)
	cm.AddContainer(kafkaContainer)
	s.Require().NoError(cm.RunAndWaitReady(s.ctx))

	s.cm = cm

	cfg.Postgres.Host = psqlContainer.GetHost()
	cfg.Postgres.Port = psqlContainer.GetPort()
	cfg.Postgres.User = psqlContainer.GetUser()
	cfg.Postgres.Password = psqlContainer.GetPassword()
	cfg.Postgres.DB = psqlContainer.GetDB()
	cfg.Postgres.SSLMode = psqlContainer.GetSSLMode()

	cfg.Kafka.Brokers = kafkaContainer.GetBrokers()

	// migrate db
	m, err := migrate.New("file://../../../migrations", psqlContainer.Dsn())
	s.Require().NoError(err)
	s.Require().NoError(m.Up())

	// aws session and s3 client
	s.downloaderAPIMock = amazon.NewMockDownloaderAPI(s.T())
	s.uploaderAPIMock = amazon.NewMockUploaderAPI(s.T())

	aws := amazon.NewAWS(&cfg.AWS)
	_ = aws.CreateSession()

	// create app
	handler, db, s3, mb, err := commands.CreateApp(cfg, log, aws, s.uploaderAPIMock, s.downloaderAPIMock)
	s.Require().NoError(err)

	s.db = db
	s.s3 = s3
	s.mb = mb

	// create token
	payload := model.TokenPayload{
		AccountID:    1,
		PatientID:    1,
		SpecialistID: 1,
	}

	accessToken, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	refreshToken, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	s.token = &model.Token{Access: *accessToken, Refresh: *refreshToken}

	// run server
	srv := httptest.NewServer(handler)
	s.client = client.NewHTTPClient(srv.URL, http.DefaultClient)
}

func (s *TestSuite) TearDownSuite() {
	s.Require().NoError(s.cm.Stop(context.Background()))
	network := s.cm.GetNetwork()
	s.Require().NoError(s.cm.RemoveNetwork(network.ID))
}
