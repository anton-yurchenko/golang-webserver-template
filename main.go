package main

import (
	"context"
	"fmt"
	"golang-webserver-template/routes"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

const (
	Version                string = "0.1.0"
	ReadTimeoutSeconds     int64  = 30
	WriteTimeoutSeconds    int64  = 30
	IdleTimeoutSeconds     int64  = 30
	ShutdownTimeoutSeconds int64  = 10
)

var (
	Port uint64 = 8080
)

func init() {
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            false,
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
	})
	logrus.SetOutput(os.Stdout)

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	if os.Getenv("PORT") != "" {
		p, err := strconv.ParseUint(os.Getenv("PORT"), 10, 64)
		if err != nil {
			logrus.Fatalf("error parsing PORT: %v", err.Error())
		}

		var maxPorts uint64 = 65535
		if p > maxPorts {
			logrus.Fatalf("port cannot be higher than %v", maxPorts)
		}

		Port = p
	}
}

func shutdown(log *logrus.Entry, srv *http.Server, timeout time.Duration) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error encountered during server shutdown: %v", err.Error())
	}

	os.Exit(0)
}

func main() {
	log := logrus.WithField("version", Version)

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	router := gin.New()
	router.GET("/_health", routes.Health)
	router.HEAD("/_health", routes.Health)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", Port),
		WriteTimeout: time.Duration(WriteTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(ReadTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(IdleTimeoutSeconds) * time.Second,
		Handler:      router,
	}

	go func() {
		log.Infof("server listening on port %v", Port)
		log.Fatal(srv.ListenAndServe())
	}()

	shutdown(log, srv, 300)
}
