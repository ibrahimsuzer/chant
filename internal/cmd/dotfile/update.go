package dotfile

import (
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/spf13/cobra"
)

type dotfileUpdateCommand struct {
	dbClient       *db.PrismaClient
	dotfileManager dotfileManager
}

func NewDotfileUpdateFactory(dbClient *db.PrismaClient, manage dotfileManager) *dotfileUpdateCommand {
	return &dotfileUpdateCommand{dbClient: dbClient, dotfileManager: manage}
}

func (f *dotfileUpdateCommand) CreateCommand() (*cobra.Command, error) {

	updateCmd := &cobra.Command{ //nolint:exhaustivestruct
		Use:   "update",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := f.dotfileManager.Update(cmd.Context(), args...)
			if err != nil {
				return fmt.Errorf("failed to update dotfiles: %w", err)
			}

			return nil
		},
	}

	return updateCmd, nil
}
