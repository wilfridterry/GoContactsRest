package config

import (
	"log"
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

	Auth Auth

	Server Server

	Grpc Grpc

	Logger Logger
}

type Auth struct {
	TokenTTL time.Duration `mapstructure:"token_ttl"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Grpc struct {
	Port int `mapstructure:"port"`
}

type Logger struct {
	Dir      string `mapstructure:"dir"`
	Filename string `mapstructure:"filename"`
}

type Postgres struct {
	Host     string
	Port     uint16
	Database string
	Username string
	Password string
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

	if err := godotenv.Load(); err == nil {
		log.Fatal("error with load env file")
	}

	viper.AutomaticEnv()

	viper.BindEnv("enviroment", "ENVIROMENT")
	viper.BindEnv("SECRET", "secret")

	viper.SetEnvPrefix("db")
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.database", "DB_DATABASE")
	viper.BindEnv("db.username", "DB_USERNAME")
	viper.BindEnv("db.password", "DB_PASSWORD")
	
	viper.SetEnvPrefix("rabbitmq")
	viper.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq.username", "RABBITMQ_USERNAME")
	viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")
	viper.BindEnv("rabbitmq.queue", "RABBITMQ_QUEUE")
	
	viper.SetEnvPrefix("server")
	viper.BindEnv("server.port", "SERVER_PORT")
	
	viper.SetEnvPrefix("grpc")
	viper.BindEnv("grpc.port", "GRPC_PORT")
	
	viper.SetEnvPrefix("auth")
	viper.BindEnv("auth.token_ttl", "AUTH_TOKEN_TTL")
	
	viper.SetEnvPrefix("logger")
	viper.BindEnv("logger.dir", "LOGGER_DIR")
	viper.BindEnv("logger.filename", "LOGGER_FILENAME")

	if err := envconfig.Process("db", &cf.DB); err != nil {
		return nil, err
	}

	if err := envconfig.Process("rabbitmq", &cf.Rabbitmq); err != nil {
		return nil, err
	}

	return cf, nil
}
