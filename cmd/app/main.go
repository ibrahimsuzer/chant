package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-enry/go-enry/v2"
	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/manage"
	"github.com/ibrahimsuzer/chant/internal/storage"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "v0.0.0"
var commit = "00000000" //nolint:gochecknoglobals

func main() {

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

	options := []ff.Option{ff.WithEnvVarPrefix("CHANT")}

	// MANAGE
	manageAddCmd := &ffcli.Command{
		Name:        "add",
		ShortUsage:  "",
		ShortHelp:   "",
		LongHelp:    "",
		UsageFunc:   nil,
		FlagSet:     nil,
		Options:     options,
		Subcommands: nil,
		Exec: func(ctx context.Context, args []string) error {
			dbClient := db.NewClient()
			if err := dbClient.Prisma.Connect(); err != nil {
				return fmt.Errorf("failed to connect to storage: %w", err)
			}

			defer func() {
				if err := dbClient.Prisma.Disconnect(); err != nil {
					panic(fmt.Errorf("failed to disconnect from storage: %w", err))
				}
			}()

			configFileRepo := storage.NewConfigFileRepo(dbClient)

			configFiles := make([]*manage.ConfigFile, 0, len(args))

			for _, path := range args {

				// Check path
				stat, err := os.Stat(path)
				if errors.Is(err, os.ErrNotExist) {
					log.Warnf("path doesn't exist: %s", path)
					continue
				} else if err != nil {
					log.Warnf("failed to read path: %s", err)
					continue
				}

				if stat.IsDir() {
					log.Warnf("cannot process directories: %s", path)
					continue
				}

				// Check details
				content, err := ioutil.ReadFile(path)
				if err != nil {
					log.Warnf("failed to read path: %s", err)
					continue
				}

				language := enry.GetLanguage(path, content)
				mimeType := enry.GetMIMEType(path, language)

				configFiles = append(configFiles, &manage.ConfigFile{
					Name:      "",
					Path:      path,
					Extension: filepath.Ext(path),
					MimeType:  mimeType,
					Language:  language,
				})
			}
			
			err := configFileRepo.Add(ctx, configFiles...)
			if err != nil {
				return fmt.Errorf("failed to write to storage: %w", err)
			}
			return nil
		},
	}

	manageCmd := &ffcli.Command{
		Name:        "manage",
		ShortUsage:  "",
		ShortHelp:   "",
		LongHelp:    "",
		UsageFunc:   nil,
		FlagSet:     nil,
		Options:     options,
		Subcommands: []*ffcli.Command{manageAddCmd},
		Exec: func(ctx context.Context, args []string) error {

			dbClient := db.NewClient()
			if err := dbClient.Prisma.Connect(); err != nil {
				return fmt.Errorf("failed to connect to storage: %w", err)
			}

			defer func() {
				if err := dbClient.Prisma.Disconnect(); err != nil {
					panic(fmt.Errorf("failed to disconnect from storage: %w", err))
				}
			}()

			configFileRepo := storage.NewConfigFileRepo(dbClient)
			list, err := configFileRepo.List(ctx, 0, 10)
			if err != nil {
				return fmt.Errorf("failed to list config files: %w", err)
			}

			fmt.Printf("%v", list)

			return nil
		},
	}

	// VERSION
	version := &ffcli.Command{
		Name:        "version",
		ShortUsage:  "",
		ShortHelp:   "",
		LongHelp:    "",
		UsageFunc:   nil,
		FlagSet:     nil,
		Options:     options,
		Subcommands: nil,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Printf("Chant %v (%v)", version, commit)
			return nil
		},
	}

	// ROOT
	rootFs := flag.NewFlagSet("chant", flag.ExitOnError)
	var (
		verbose = rootFs.Bool("verbose", false, "increase log verbosity")
	)

	root := &ffcli.Command{
		Name:        "chant",
		ShortUsage:  "chant [-version] [-help] [-autocomplete-(un)install] <command> [args]",
		ShortHelp:   "Chant: Local environment manager for dotfiles and project binaries",
		LongHelp:    "",
		UsageFunc:   nil,
		FlagSet:     rootFs,
		Options:     options,
		Subcommands: []*ffcli.Command{version, manageCmd},
		Exec: func(ctx context.Context, args []string) error {
			fmt.Printf("verbose %v\n", *verbose)
			return nil
		},
	}

	// RUN
	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
