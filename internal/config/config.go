package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Enviroment string
	Secret     string

	DB Postgres

	Rabbitmq Rabbitmq

	Auth struct {
		TokenTTL time.Duration `mapstructure:"ttl"`
	} `mapstructure:"auth"`

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Grpc struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"grpc"`

	Logger struct {
		Dir      string `mapstructure:"dir"`
		Filename string `mapstructure:"filename"`
	} `mapstructure:"logger"`
}

type Postgres struct {
	Host     string
	Port     uint16
	Database string
	Username string
	Password string
	SSLMode  bool
}

type Rabbitmq struct {
	Host     string
	Port     uint16
	Queue    string
	Username string
	Password string
}

func NewConfig(dirname, filename string) (*Config, error) {
	cf := new(Config)
	viper.AddConfigPath(dirname)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cf); err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := envconfig.Process("db", &cf.DB); err != nil {
		return nil, err
	}

	if err := envconfig.Process("rabbitmq", &cf.Rabbitmq); err != nil {
		return nil, err
	}

	return cf, nil
}
