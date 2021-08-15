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

func (s *fileVersion) Add(ctx context.Context, fileId string, content string) (*dotfiles.FileVersion, error) {
	hash := hash256(content)

	createFileVersion, err := s.dbClient.FileVersion.CreateOne(
		db.FileVersion.Content.Set(content),
		db.FileVersion.Hash.Set(hash),
		db.FileVersion.File.Link(db.Dotfile.ID.Equals(fileId)),
		db.FileVersion.CurrentlyUsed.Link(db.Dotfile.ID.Equals(fileId)),
	).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), uniqueConstraintViolation) {
			return nil, storage_errors.ErrUniqueConstraintViolation
		}

		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return mapModelToFileVersion(createFileVersion), nil

}

func (s *fileVersion) Update(ctx context.Context, fileId string, content string) (*dotfiles.FileVersion, error) {
	hash := hash256(content)

	dotfile, err := s.dbClient.Dotfile.FindUnique(db.Dotfile.ID.Equals(fileId)).With(db.Dotfile.Current.Fetch()).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot find the file: %w", err)
	}

	// If current exists, check if hash is same and return
	var foundCurrent *string
	current, ok := dotfile.Current()
	if ok && current.Hash == hash {
		return mapModelToFileVersion(current), nil
	} else if ok {
		foundCurrent = &current.ID // If current exists, use it as predecessor
	}

	// If current doesn't exist or hash is different
	// add new version and update current

	createFileVersion, err := s.dbClient.FileVersion.CreateOne(
		db.FileVersion.Content.Set(content),
		db.FileVersion.Hash.Set(hash),
		db.FileVersion.File.Link(db.Dotfile.ID.Equals(fileId)),
		db.FileVersion.CurrentlyUsed.Link(db.Dotfile.ID.Equals(fileId)),
		db.FileVersion.Predecessor.Link(db.FileVersion.ID.EqualsIfPresent(foundCurrent)),
	).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to add new version: %w", err)
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
