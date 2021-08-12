package dotfile

import (
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/spf13/cobra"
)

type dotfileRemoveCommand struct {
	dbClient       *db.PrismaClient
	dotfileManager dotfileManager
}

func NewDotfileRemoveFactory(dbClient *db.PrismaClient, manage dotfileManager) *dotfileRemoveCommand {
	return &dotfileRemoveCommand{dbClient: dbClient, dotfileManager: manage}
}

func (f *dotfileRemoveCommand) CreateCommand() (*cobra.Command, error) {
	dotfileListCmd := &cobra.Command{ //nolint:exhaustivestruct
		Use:   "remove",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := f.dotfileManager.Remove(cmd.Context(), args...)
			if err != nil {
				return fmt.Errorf("failed to remove dotfiles: %w", err)
			}

			return nil
		},
	}

	return dotfileListCmd, nil
}
