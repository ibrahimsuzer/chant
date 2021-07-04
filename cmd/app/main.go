package main

import (
	"os"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/cmd"
	"github.com/ibrahimsuzer/chant/internal/cmd/manage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "v0.0.0"
var commit = "" //nolint:gochecknoglobals

func main() {

	// TODO: Replace logger with colored output
	// TODO: Add structure for managinc config files and creator method for cobra command

	// printer := pterm.BasicTextPrinter{}

	// Logger
	logEncoder := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "T",
		NameKey:        "logger",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	baseLogger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(logEncoder), os.Stdout, zap.InfoLevel))
	defer baseLogger.Sync() // flushes buffer, if any
	log := baseLogger.Sugar()

	// Initialize Root Command and Configuration
	v := viper.New()
	cfg := cmd.NewConfiguration(v)
	rootFactory := cmd.NewRootFactory(log, cfg, version, commit)
	rootCmd, err := rootFactory.CreateCommand()
	if err != nil {
		log.Fatalf("failed to create command: %v", err)
	}

	// Manage Command
	dbClient := db.NewClient()
	manageFactory := manage.NewManageFactory(dbClient)
	manageCmd, err := manageFactory.CreateCommand()
	if err != nil {
		log.Fatalf("failed to create command: %v", err)
	}

	// Manage Add
	manageListFactory := manage.NewManageAddFactory(dbClient)
	manageListCmd, err := manageListFactory.CreateCommand()
	if err != nil {
		log.Fatalf("failed to create command: %v", err)
	}

	manageCmd.AddCommand(manageListCmd)
	rootCmd.AddCommand(manageCmd)
	err = rootCmd.Execute()
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}

}
