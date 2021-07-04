package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type logger interface {
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type rootCommand struct {
	version, commit string
	cfg             *configuration
	log             logger
}

func NewRootFactory(log logger, cfg *configuration, version string, commit string) *rootCommand {
	return &rootCommand{version: version, commit: commit, cfg: cfg, log: log}
}

func (f *rootCommand) CreateCommand() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "chant",
		Short: "Chant: Local environment manager for dotfiles and project binaries",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			f.log.Infof("verbose enabled %t", f.cfg.App.Verbose)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			f.cfg.Configure()
			err := f.cfg.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %v", err)
			}

			err = presetRequiredFlags(cmd)
			if err != nil {
				return fmt.Errorf("failed update prequired flags: %v", err)
			}
			return nil
		},
	}
	rootCmd.Version = fmt.Sprintf("%v (%v)", f.version, f.commit)
	rootCmd.PersistentFlags().BoolVarP(&f.cfg.App.Verbose, "verbose", "v", false, "verbose output")
	f.cfg.viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	return rootCmd, nil

}

func presetRequiredFlags(cmd *cobra.Command) error {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return fmt.Errorf("error binding flags: %v", err)
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			err = cmd.Flags().Set(f.Name, viper.GetString(f.Name))
			if err != nil {
				log.Fatalf("error setting flag %s: %v", f.Name, err)
			}
		}
	})

	return nil
}
