package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	Server   Server
	Postgres Postgres
	Logger   Logger
	AWS      AWS
	JWT      JWT
	Verify   Verify
	Kafka    Kafka
}

type Server struct {
	AppVersion        string
	Host              string
	Port              string
	Development       bool
	Debug             bool
	SSL               bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CtxDefaultTimeout time.Duration
	IdleTimeout       time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

type Postgres struct {
	Host            string
	Port            string
	User            string
	Password        string
	DB              string
	SSLMode         string
	Driver          string
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

type Logger struct {
	LoggerName        string
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type AWS struct {
	Region            string
	ID                string
	SecretAccessKey   string
	AccountBucketName string
}

type JWT struct {
	AccessTokenCookieName  string
	AccessTokenSecretKey   string
	AccessTokenExpiresAt   time.Duration
	RefreshTokenCookieName string
	RefreshTokenSecretKey  string
	RefreshTokenExpiresAt  time.Duration
}

type Verify struct {
	VerifyCodeCookieName string
	VerifyCodeExpiresAt  time.Duration
}

type Kafka struct {
	Brokers  []string
	Deadline time.Duration
}

func Get(filename string, path ...string) (*Config, error) {
	if len(path) == 0 {
		viper.AddConfigPath("./config")
	} else {
		viper.AddConfigPath(strings.Join(path, ""))
	}

	viper.SetConfigName(filename)

	viper.AutomaticEnv()                                   // get vars from environment and replace values in viper structure if the keys match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // help viper define keys

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
