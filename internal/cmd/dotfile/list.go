package dotfile

import (
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/spf13/cobra"
)

type dotfileListCommand struct {
	dbClient       *db.PrismaClient
	dotfileManager dotfileManager
}

func NewDotfileListFactory(dbClient *db.PrismaClient, manage dotfileManager) *dotfileListCommand {
	return &dotfileListCommand{dbClient: dbClient, dotfileManager: manage}
}

func (f *dotfileListCommand) CreateCommand() (*cobra.Command, error) {
	dotfileCmd := &cobra.Command{
		Use:     "list",
		Short:   "",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := f.dotfileManager.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to list dotfiles: %w", err)
			}

			fmt.Printf("%v", list)
			return nil
		},
	}
	return dotfileCmd, nil
}
