package manage

import (
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/storage"
	"github.com/spf13/cobra"
)

type manageCommand struct {
	dbClient *db.PrismaClient
}

func NewManageFactory(dbClient *db.PrismaClient) *manageCommand {
	return &manageCommand{dbClient: dbClient}
}

func (f *manageCommand) CreateCommand() (*cobra.Command, error) {
	manageCmd := &cobra.Command{
		Use:     "manage",
		Short:   "",
		Aliases: []string{"m"},

		RunE: func(cmd *cobra.Command, args []string) error {
			configFileRepo := storage.NewConfigFileRepo(f.dbClient)
			list, err := configFileRepo.List(cmd.Context(), 0, 10)
			if err != nil {
				return fmt.Errorf("failed to list config files: %w", err)
			}

			fmt.Printf("%v", list)
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
	return manageCmd, nil
}
