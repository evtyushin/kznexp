package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"kznexp/api/apiserver"
	"kznexp/service"
)

const (
	loggerPrdLvl = "prod"
	loggerDevLvl = "dev"
)

func main() {

	configPath := flag.String("config", "configs/config.yaml", "Config path")
	flag.Parse()

	config, err := readConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read config %e", err))
	}

	logger, err := newLogger()(config)
	if err != nil {
		panic(fmt.Sprintf("failed to get logger %e", err))
	}

	serv := service.New(config.Count, config.IntervalSec)
	chanTasks := serv.CreateChanTasks()
	http := apiserver.New(serv, config.HttpServer.Listen, chanTasks, logger)

	go func() {
		if err := http.Start(); err != nil {
			logger.Error("failed to continue serving", zap.Error(err))
		}
	}()

	go func() {
		serv.StartProcessing(chanTasks)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("service shutdown")
}

type conf struct {
	HttpServer struct {
		Listen string `yaml:"listen"`
	} `yaml:"http_server"`
	Count        int    `yaml:"count"`
	IntervalSec  int    `yaml:"interval_sec"`
	LoggingLevel string `yaml:"logging_level"`
}

func readConfig(path string) (*conf, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(abs)
	if err != nil {
		return nil, err
	}

	var config *conf
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func newLogger() func(config *conf) (*zap.Logger, error) {
	return func(config *conf) (*zap.Logger, error) {
		loggerConfig := zap.NewDevelopmentConfig()

		switch config.LoggingLevel {
		case loggerPrdLvl:
			loggerConfig = zap.NewProductionConfig()
		case loggerDevLvl:
			loggerConfig = zap.NewDevelopmentConfig()
		default:
			return nil, fmt.Errorf("unknown logger setup %s", config.LoggingLevel)
		}

		l, err := loggerConfig.Build()
		if err != nil {
			return nil, fmt.Errorf("failed to build logger config: %w", err)
		}
		return l, nil
	}
}
