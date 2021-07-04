package dotfile

import (
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/spf13/cobra"
)

type dotfileAddCommand struct {
	dbClient       *db.PrismaClient
	dotfileManager dotfileManager
}

func NewDotfileAddFactory(dbClient *db.PrismaClient, manage dotfileManager) *dotfileAddCommand {
	return &dotfileAddCommand{dbClient: dbClient, dotfileManager: manage}
}

func (f *dotfileAddCommand) CreateCommand() (*cobra.Command, error) {
	dotfileListCmd := &cobra.Command{
		Use:   "add",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := f.dotfileManager.Add(cmd.Context(), args...)
			if err != nil {
				return fmt.Errorf("failed to list dotfiles: %w", err)
			}

			return nil
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := f.dbClient.Prisma.Connect(); err != nil {
				return fmt.Errorf("failed to connect to storage: %w", err)
			}
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if err := f.dbClient.Prisma.Disconnect(); err != nil {
				return fmt.Errorf("failed to disconnect from storage: %w", err)
			}
			return nil
		},
	}
	return dotfileListCmd, nil
}
