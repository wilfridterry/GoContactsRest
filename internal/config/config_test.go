package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	type env struct {
		enviroment       string
		secret           string
		dbHost           string
		dbPort           string
		dbDatabase       string
		dbUsername       string
		dbPassword       string
		rabbitmqHost     string
		rabbitmqPort     string
		rabbitmqQueue    string
		rabbitmqUsername string
		rabbitmqPassword string
		authTokenTTL     string
		serverPort       string
		grpcPort         string
		loggerDir        string
		loggerFilename   string
	}

	type args struct {
		env  env
		path string
		filename string
	}

	setEnv := func(env env) {
		if env.enviroment != "" {
			os.Setenv("ENVIROMENT", env.enviroment)
		}
		if env.secret != "" {
			os.Setenv("SECRET", env.secret)
		}
		if env.dbHost != "" {
			os.Setenv("DB_HOST", env.dbHost)
		}
		if env.dbPort != "" {
			os.Setenv("DB_PORT", env.dbPort)
		}
		if env.dbDatabase != "" {
			os.Setenv("DB_DATABASE", env.dbDatabase)
		}
		if env.dbUsername != "" {
			os.Setenv("DB_USERNAME", env.dbUsername)
		}
		if env.dbPassword != "" {
			os.Setenv("DB_PASSWORD", env.dbPassword)
		}
		if env.rabbitmqHost != "" {
			os.Setenv("RABBITMQ_HOST", env.rabbitmqHost)
		}
		if env.rabbitmqPort != "" {
			os.Setenv("RABBITMQ_PORT", env.rabbitmqPort)
		}
		if env.rabbitmqQueue != "" {
			os.Setenv("RABBITMQ_QUEUE", env.rabbitmqQueue)
		}
		if env.rabbitmqUsername != "" {
			os.Setenv("RABBITMQ_USERNAME", env.rabbitmqUsername)
		}
		if env.rabbitmqPassword != "" {
			os.Setenv("RABBITMQ_PASSWORD", env.rabbitmqPassword)
		}
		if env.authTokenTTL != "" {
			os.Setenv("AUTH_TOKEN_TTL", env.authTokenTTL)
		}
		if env.serverPort != "" {
			os.Setenv("SERVER_PORT", env.serverPort)
		}
		if env.grpcPort != "" {
			os.Setenv("GRPC_PORT", env.grpcPort)
		}
		if env.loggerDir != "" {
			os.Setenv("LOGGER_DIR", env.loggerDir)
		}
		if env.loggerFilename != "" {
			os.Setenv("LOGGER_FILENAME", env.loggerFilename)
		}
	}

	testCases := []struct{
		name string
		args args
		want *Config
		wantErr bool
	}{
		{
			name: "test from config file",
			args: args{
				path: "fixtures",
				filename: "main",
			},
			want: &Config{
				Enviroment: "testing",
				Secret: "salt",
				DB: Postgres{
					Host: "localhost",
					Port: 5432,
					Database: "postgres",
					Username: "root",
					Password: "password",
				},
				Rabbitmq: Rabbitmq{
					Host: "localhost",
					Port: 5672,
					Queue: "queue",
					Username: "root",
					Password: "password",
				},
				Auth: Auth{
					TokenTTL: time.Minute * 15,
				},
				Server: Server{
					Port: 8081,
				},
				Grpc: Grpc{
					Port: 9000,
				},
				Logger: Logger{
					Dir: "storage/logs",
					Filename: "test.log",
				},
			},
			wantErr: false,
		},
		{
			name: "test from env file",
			args: args{
				path: "fixtures",
				filename: "main",
				env: env{
					enviroment: "env_testing",
					secret: "env_salt",
					dbHost: "127.0.0.1",
					dbPort: "5435",
					dbDatabase: "env_postgres",
					dbUsername: "env_root",
					dbPassword: "env_password",
					rabbitmqHost: "127.0.0.1",
					rabbitmqPort: "5675",
					rabbitmqQueue: "env_queue",
					rabbitmqUsername: "env_root",
					rabbitmqPassword: "env_password",
					authTokenTTL: "30m",
					serverPort: "8082",
					grpcPort: "9001",
					loggerDir: "storage/env_logs",
					loggerFilename: "env_test.log",
				},
			},
			want: &Config{
				Enviroment: "env_testing",
				Secret: "env_salt",
				DB: Postgres{
					Host: "127.0.0.1",
					Port: 5435,
					Database: "env_postgres",
					Username: "env_root",
					Password: "env_password",
				},
				Rabbitmq: Rabbitmq{
					Host: "127.0.0.1",
					Port: 5675,
					Queue: "env_queue",
					Username: "env_root",
					Password: "env_password",
				},
				Auth: Auth{
					TokenTTL: time.Minute * 30,
				},
				Server: Server{
					Port: 8082,
				},
				Grpc: Grpc{
					Port: 9001,
				},
				Logger: Logger{
					Dir: "storage/env_logs",
					Filename: "env_test.log",
				},
			},
			wantErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setEnv(testCase.args.env)
			
			got, err := NewConfig(testCase.args.path, testCase.args.filename)
			if (err != nil) != testCase.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, testCase.wantErr)
				return 
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("NewConfig() got = %v, want = %v", got, testCase.want)
			}
		})
	}
}
