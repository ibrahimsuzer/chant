package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ibrahimsuzer/chant/db"
	"github.com/ibrahimsuzer/chant/internal/dotfiles"
	storage_errors "github.com/ibrahimsuzer/chant/internal/storage/errors"
)

type fileVersion struct {
	dbClient *db.PrismaClient
}

func NewFileVersionRepo(db *db.PrismaClient) *fileVersion {
	return &fileVersion{dbClient: db}
}

func (s *fileVersion) Add(ctx context.Context, fileId string, file *dotfiles.FileVersion) (*dotfiles.FileVersion, error) {
	hash := hash256(file.Content)

	createFileVersion, err := s.dbClient.FileVersion.CreateOne(
		db.FileVersion.Content.Set(file.Content),
		db.FileVersion.Hash.Set(hash),
		db.FileVersion.File.Link(db.Dotfile.ID.Equals(fileId)),
	).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), uniqueConstraintViolation) {
			return nil, storage_errors.ErrUniqueConstraintViolation
		}

		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return mapModelToFileVersion(createFileVersion), nil

}

func mapModelToFileVersion(model *db.FileVersionModel) *dotfiles.FileVersion {
	return &dotfiles.FileVersion{
		Id:      model.ID,
		Content: model.Content,
		Hash:    model.Hash,
	}
}

func hash256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
