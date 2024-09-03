package app

import (
	"contact-list/internal/config"
	"contact-list/internal/repository/psql"
	"contact-list/internal/service"
	grpc_client "contact-list/internal/transport/grpc"
	"contact-list/internal/transport/rest"
	amqplog "contact-list/pkg/amqp_log"
	"contact-list/pkg/database"
	"contact-list/pkg/hashier"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	CONFDIR      = "configs"
	CONFFILENAME = "main"
)

func initLogger(dir, filename string) (io.Closer, error) {
	_ = os.Mkdir(dir, 0666)
	path := filepath.Join(dir, filename)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(file)
	log.SetLevel(log.DebugLevel)

	return file, nil
}

func Run() {
	ctx := context.Background()

	cf, err := config.NewConfig(CONFDIR, CONFFILENAME)
	if err != nil {
		log.Error(err)
	}

	filelog, err := initLogger(cf.Logger.Dir, cf.Logger.Filename)
	if err != nil {
		log.Error(err)
	}
	defer filelog.Close()

	conn, err := database.NewConnection(ctx, &database.ConnectionConfig{
		Host:     cf.DB.Host,
		Port:     cf.DB.Port,
		Database: cf.DB.Database,
		Username: cf.DB.Username,
		Password: cf.DB.Password,
	})

	if err != nil {
		log.Error(err)
	}
	defer conn.Close(ctx)

	auditClient, err := grpc_client.NewClient(cf.Grpc.Port)
	if err != nil {
		log.Error(err)
	}

	amqpClient, err := amqplog.New(&amqplog.ConfigOptions{
		Host: cf.Rabbitmq.Host,
		Port: int(cf.Rabbitmq.Port),
		Username: cf.Rabbitmq.Username,
		Password: cf.Rabbitmq.Password,
		Queue: cf.Rabbitmq.Queue,
	})
	if err != nil {
		log.Error(err)
	}

	auditLogService := service.NewAuditLog(amqpClient)
	contactsRepo := psql.NewContacts(conn)
	contactsService := service.NewContacts(contactsRepo, auditClient, auditLogService)

	userRepo := psql.NewUsers(conn)
	hashier := hashier.NewHashier(cf.Secret)
	sessionRepo := psql.NewTokens(conn)
	authService := service.New(userRepo, sessionRepo, auditClient, auditLogService, hashier, []byte(cf.Secret), cf.Auth.TokenTTL)

	handler := rest.NewHandler(contactsService, authService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cf.Server.Port),
		Handler: handler.InitRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.WithField("error", err).Fatal("listening err")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithField("error", err).Fatal("Server forced to shutdown:")
	}

	log.Info("Exiting server")
}
