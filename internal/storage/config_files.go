package storage

import (
	"context"
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/manage"
	"github.com/prisma/prisma-client-go/runtime/transaction"
)

type configFileRepo struct {
	dbClient *db.PrismaClient
}

func NewConfigFileRepo(db *db.PrismaClient) *configFileRepo {
	return &configFileRepo{dbClient: db}
}

func (s *configFileRepo) Add(ctx context.Context, files ...*manage.ConfigFile) error {

	queries := make([]transaction.Param, 0, len(files))

	for _, file := range files {

		createConfigFile := s.dbClient.ConfigFile.CreateOne(
			db.ConfigFile.Name.Set(file.Name),
			db.ConfigFile.Description.Set(file.Description),
			db.ConfigFile.Location.Set(file.Location),
			db.ConfigFile.Format.Link(db.Format.Name.Equals(string(file.Format))),
		).Tx()

		queries = append(queries, createConfigFile)
	}

	err := s.dbClient.Prisma.Transaction(queries...).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *configFileRepo) List(ctx context.Context, page, count int) ([]*manage.ConfigFile, error) {

	configFileModels, err := s.dbClient.ConfigFile.FindMany().Take(count).Skip(page * count).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	result := make([]*manage.ConfigFile, 0, count)

	for _, configFileModel := range configFileModels {

		configFile := manage.ConfigFile{
			Id:          configFileModel.ID,
			Name:        configFileModel.Name,
			Description: configFileModel.Description,
			Location:    configFileModel.Location,
			Format:      manage.FileFormat(configFileModel.Format().Name),
		}

		result = append(result, &configFile)
	}

	return result, nil

}
