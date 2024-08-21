package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Enviroment string
	Secret string

	DB Postgres

	Auth struct {
		TokenTTL time.Duration `mapstructure:"ttl"`
	} `mapstructure:"auth"`

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Logger struct {
		Dir string `mapstructure:"dir"`
		Filename string `mapstructure:"filename"`
	} `mapstructure:"logger"`
}

type Postgres struct {
	Host     string
	Port     uint16
	Database   string
	Username string
	Password string
	SSLMode  bool 
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

	return cf, nil
}