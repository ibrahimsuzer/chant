package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/genjidb/genji/sql/driver"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var version = "v0.0.0"
var commit = "00000000" //nolint:gochecknoglobals

func main() {

	options := []ff.Option{ff.WithEnvVarPrefix("CHANT")}

	// Create a sql/database DB instance
	db, err := sql.Open("genji", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	

	// MANAGE
	add := &ffcli.Command{
		Name:        "manage",
		ShortUsage:  "",
		ShortHelp:   "",
		LongHelp:    "",
		UsageFunc:   nil,
		FlagSet:     nil,
		Options:     options,
		Subcommands: nil,
		Exec: func(ctx context.Context, args []string) error {
			
			
			
			
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
		Subcommands: []*ffcli.Command{version},
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
