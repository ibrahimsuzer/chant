package dotfile

import (
	"context"
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/dotfiles"
	"github.com/spf13/cobra"
)

type dotfileManager interface {
	Add(ctx context.Context, files ...string) error
	List(ctx context.Context) ([]*dotfiles.Dotfile, error)
}

type dotfileCommand struct {
	dbClient       *db.PrismaClient
	dotfileManager dotfileManager
}

func NewDotfileCommandFactory(dbClient *db.PrismaClient, manage dotfileManager) *dotfileCommand {
	return &dotfileCommand{dbClient: dbClient, dotfileManager: manage}
}

func (f *dotfileCommand) CreateCommand() (*cobra.Command, error) {
	dotfileCmd := &cobra.Command{
		Use:     "dotfile",
		Short:   "",
		Aliases: []string{"d"},
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
	return dotfileCmd, nil
}
