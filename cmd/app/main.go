package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/cmd"
	"github.com/ibrahimsuzer/chant/internal/cmd/dotfile"
	"github.com/ibrahimsuzer/chant/internal/dotfiles"
	"github.com/ibrahimsuzer/chant/internal/printer"
	"github.com/spf13/viper"
)

var version = "v0.0.0"
var commit = "" //nolint:gochecknoglobals

func main() {

	// TODO: Replace logger with colored output
	// TODO: Add structure for managinc config files and creator method for cobra command

	// printer := pterm.BasicTextPrinter{}

	// Initialize Root Command and Configuration
	v := viper.New()
	cfg := cmd.NewConfiguration(v)
	rootFactory := cmd.NewRootFactory(cfg, version, commit)
	rootCmd, err := rootFactory.CreateCommand()
	if err != nil {
		log.Fatalf("failed to create command: %v", err)
	}

	// Manage Command
	dbClient := db.NewClient()
	dotfileRepo := dotfiles.NewDotFileRepo(dbClient)
	dotfilePrinter := printer.NewPrinter(color.New(color.Reset))
	dotfileManager := dotfiles.NewDotfileManager(dotfileRepo, dotfilePrinter)
	dotfileCommandFactory := dotfile.NewDotfileCommandFactory(dbClient, dotfileManager)
	dotfileCmd, err := dotfileCommandFactory.CreateCommand()
	if err != nil {
		log.Fatalf("failed to create command: %v", err)
	}

	// Manage Add
	dotfileAddFactory := dotfile.NewDotfileAddFactory(dbClient, dotfileManager)
	dotfileAddCmd, err := dotfileAddFactory.CreateCommand()
	if err != nil {
		log.Fatalf("failed to create command: %v", err)
	}

	// Manage Add
	dotfileListFactory := dotfile.NewDotfileListFactory(dbClient, dotfileManager)
	dotfileListCmd, err := dotfileListFactory.CreateCommand()
	if err != nil {
		log.Fatalf("failed to create command: %v", err)
	}

	dotfileCmd.AddCommand(dotfileAddCmd)
	dotfileCmd.AddCommand(dotfileListCmd)
	rootCmd.AddCommand(dotfileCmd)
	err = rootCmd.Execute()
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}

}
