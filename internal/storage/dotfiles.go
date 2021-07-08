package storage

import (
	"context"
	"fmt"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/dotfiles"
	"github.com/prisma/prisma-client-go/runtime/transaction"
)

type dotfileRepo struct {
	dbClient *db.PrismaClient
}

func NewDotFileRepo(db *db.PrismaClient) *dotfileRepo {
	return &dotfileRepo{dbClient: db}
}

func (s *dotfileRepo) Add(ctx context.Context, files ...*dotfiles.Dotfile) error {

	queries := make([]transaction.Param, 0, len(files))

	for _, file := range files {

		createDotfile := s.dbClient.Dotfile.CreateOne(
			db.Dotfile.Name.Set(file.Name),
			db.Dotfile.Path.Set(file.Path),
			db.Dotfile.Extension.Set(file.Extension),
			db.Dotfile.MimeType.Set(file.MimeType),
			db.Dotfile.Language.Set(file.Language),

		).Tx()

		queries = append(queries, createDotfile)
	}

	err := s.dbClient.Prisma.Transaction(queries...).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *dotfileRepo) List(ctx context.Context, page, count int) ([]*dotfiles.Dotfile, error) {

	// Base query
	query := s.dbClient.Dotfile.FindMany()

	// Pagination
	if count > 0 {
		query = query.Take(count).Skip(page * count)
	}

	dotfileModels, err := query.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	result := make([]*dotfiles.Dotfile, 0, count)
	for _, dotfileModel := range dotfileModels {

		dotfile := dotfiles.Dotfile{
			Id:        dotfileModel.ID,
			Name:      dotfileModel.Name,
			Path:      dotfileModel.Path,
			Extension: dotfileModel.Extension,
			MimeType:  dotfileModel.MimeType,
			Language:  dotfileModel.Language,
		}

		result = append(result, &dotfile)
	}

	return result, nil

}
