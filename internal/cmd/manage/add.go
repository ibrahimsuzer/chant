package manage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-enry/go-enry/v2"
	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/manage"
	"github.com/ibrahimsuzer/chant/internal/storage"
	"github.com/spf13/cobra"
)

type manageAddCommand struct {
	dbClient *db.PrismaClient
}

func NewManageAddFactory(dbClient *db.PrismaClient) *manageAddCommand {
	return &manageAddCommand{dbClient: dbClient}
}

func (f *manageAddCommand) CreateCommand() (*cobra.Command, error) {
	manageListCmd := &cobra.Command{
		Use:   "add",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := cmd.Context()

			configFileRepo := storage.NewConfigFileRepo(f.dbClient)

			configFiles := make([]*manage.ConfigFile, 0, len(args))

			for _, path := range args {

				// Check path
				stat, err := os.Stat(path)
				if errors.Is(err, os.ErrNotExist) {
					fmt.Printf("path doesn't exist: %s", path)
					continue
				} else if err != nil {
					fmt.Printf("failed to read path: %s", err)
					continue
				}

				if stat.IsDir() {
					fmt.Printf("cannot process directories: %s", path)
					continue
				}

				// Check details
				content, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Printf("failed to read path: %s", err)
					continue
				}

				language := enry.GetLanguage(path, content)
				mimeType := enry.GetMIMEType(path, language)

				configFiles = append(configFiles, &manage.ConfigFile{
					Name:      "",
					Path:      path,
					Extension: filepath.Ext(path),
					MimeType:  mimeType,
					Language:  language,
				})
			}

			err := configFileRepo.Add(ctx, configFiles...)
			if err != nil {
				return fmt.Errorf("failed to write to storage: %w", err)
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
	return manageListCmd, nil
}
