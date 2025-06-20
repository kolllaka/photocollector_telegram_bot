package main

import (
	"fmt"
	"io"
	"os"
	"time"

	tgClient "github.com/kolllaka/telegram_bot/internal/clients/telegram"
	eventConsumer "github.com/kolllaka/telegram_bot/internal/consumer/event-consumer"
	"github.com/kolllaka/telegram_bot/internal/events/telegram"
	"github.com/kolllaka/telegram_bot/internal/model"
	tgStorage "github.com/kolllaka/telegram_bot/internal/storage/telegram"
	"github.com/kolllaka/telegram_bot/pkg/config"
	"github.com/kolllaka/telegram_bot/pkg/logging"
)

const (
	batchSize     = 10
	YAML_CFG_PATH = "./locale.yaml"
	LOGS_FOLDER   = "./logs"
)

func main() {
	// env config init
	envCfg := model.NewEnvConfig()
	if err := config.LoadEnv(envCfg); err != nil {
		os.Exit(1)
	}

	// logger init
	fileLogger := logging.NewFileWriter(
		fmt.Sprintf("%s/%s.log", LOGS_FOLDER, time.Now().Format("20060102T1504")),
	)
	multiWriter := io.MultiWriter(fileLogger, os.Stdout)
	logger := logging.NewLogger(
		multiWriter,
		logging.WithLevel(envCfg.LogLvl),
		logging.WithIsJSON(true),
		logging.WithAddSource(false),
	)

	logger.Info("start application")

	// locale config init
	localeCfg := model.NewLocaleConfig()
	if err := config.LoadYamlByPath(YAML_CFG_PATH, localeCfg); err != nil {
		logger.Error("failed to load locale config", logging.AnyAttr("localeCfg", localeCfg))

		os.Exit(1)
	}

	// tgClient init
	tgClient := tgClient.New(envCfg.Token)

	// storage
	tgStorage := tgStorage.New(envCfg.ChannelId, tgClient, localeCfg.Infos)

	eventsProcessor := telegram.New(logger, tgClient, tgStorage, localeCfg.Commands)

	if eventsProcessor != nil {
		logger.Debug("eventsProcessor not nil")
	}

	consumer := eventConsumer.New(logger, eventsProcessor, eventsProcessor, batchSize)

	if consumer != nil {
		logger.Debug("consumer not nil")
	}

	consumer.Start()
}
