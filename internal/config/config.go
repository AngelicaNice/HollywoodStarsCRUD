package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Postgres struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

type Config struct {
	DB     Postgres
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Auth struct {
		TokenTtl time.Duration `mapstructure:"token_ttl"`
	} `mapstructure:"auth"`
}

func NewConfig(folder string, filename string) (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		log.WithField(".env", "wrong environment variables").Fatal(err)

		return nil, err
	}

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.WithField("config", "wrong config").Fatal(err)

		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		log.WithField("config", "wrong unmarshalling").Fatal(err)

		return nil, err
	}

	return cfg, nil
}
