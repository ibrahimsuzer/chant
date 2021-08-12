package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/dotfiles"
	storage_errors "github.com/ibrahimsuzer/chant/internal/storage/errors"
	"github.com/prisma/prisma-client-go/runtime/transaction"
)

const (
	uniqueConstraintViolation = "UniqueConstraintViolation"
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
		if strings.Contains(err.Error(), uniqueConstraintViolation) {
			return storage_errors.ErrUniqueConstraintViolation
		}

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

	result := mapModelsToDotfile(dotfileModels)

	return result, nil

}

func (s *dotfileRepo) Remove(ctx context.Context, ids ...string) error {

	_, err := s.dbClient.Dotfile.FindMany(db.Dotfile.ID.In(ids)).Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return nil
}

func (s *dotfileRepo) Find(ctx context.Context, ids ...string) ([]*dotfiles.Dotfile, error) {

	found, err := s.dbClient.Dotfile.FindMany(db.Dotfile.ID.In(ids)).Exec(ctx)
	if err != nil {
		return []*dotfiles.Dotfile{}, fmt.Errorf("query failed: %w", err)
	}

	result := mapModelsToDotfile(found)

	return result, nil
}

func mapModelsToDotfile(dotfileModels []db.DotfileModel) []*dotfiles.Dotfile {
	result := make([]*dotfiles.Dotfile, 0, len(dotfileModels))
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
	return result
}
