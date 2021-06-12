package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var version = "v0.0.0"
var commit = "00000000" //nolint:gochecknoglobals

func main() {
	// VERSION
	version := &ffcli.Command{
		Name:        "version",
		ShortUsage:  "",
		ShortHelp:   "",
		LongHelp:    "",
		UsageFunc:   nil,
		FlagSet:     nil,
		Options:     []ff.Option{ff.WithEnvVarPrefix("CHANT")},
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
		Options:     []ff.Option{ff.WithEnvVarPrefix("CHANT")},
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
