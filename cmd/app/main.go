package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/storage"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var version = "v0.0.0"
var commit = "00000000" //nolint:gochecknoglobals

func main() {

	options := []ff.Option{ff.WithEnvVarPrefix("CHANT")}

	// MANAGE

	manage := &ffcli.Command{
		Name:        "manage",
		ShortUsage:  "",
		ShortHelp:   "",
		LongHelp:    "",
		UsageFunc:   nil,
		FlagSet:     nil,
		Options:     options,
		Subcommands: []*ffcli.Command{},
		Exec: func(ctx context.Context, args []string) error {

			dbClient := db.NewClient()
			if err := dbClient.Prisma.Connect(); err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err := dbClient.Prisma.Disconnect(); err != nil {
					panic(err)
				}
			}()
			
			configFileRepo := storage.NewConfigFileRepo(dbClient)
			list, err := configFileRepo.List(ctx, 0, 10)
			if err != nil {
				return err
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
		Subcommands: []*ffcli.Command{version, manage},
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
