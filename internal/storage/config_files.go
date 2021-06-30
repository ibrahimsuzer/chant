package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/genjidb/genji/sql/driver"
	"github.com/ibrahimsuzer/chant/internal/manage"
)


type configFileRepo struct {
	db          *sql.DB
	idGenerator *ulidGenerator
}

func NewConfigFileRepo(db *sql.DB, idGenerator *ulidGenerator) *configFileRepo {
	return &configFileRepo{db: db, idGenerator: idGenerator}
}

func (s *configFileRepo) Add(ctx context.Context, files ...*manage.ConfigFile) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer tx.Rollback()

	for _, file := range files {
		id, err := s.idGenerator.Generate()
		if err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}

		configFile := NewConfigFile(file, id)
		_, err = tx.ExecContext(ctx, `INSERT INTO config_files VALUES ?`, &configFile)
		if err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *configFileRepo) List(ctx context.Context) ([]*manage.ConfigFile, error) {

	stream, err := s.db.QueryContext(ctx, "SELECT * FROM config_files")
	if err != nil {
		return []*manage.ConfigFile{}, fmt.Errorf("failed to list configs: %w", err)

	}
	defer stream.Close()

	var files []*manage.ConfigFile

	for stream.Next() {
		var u configFile

		err := stream.Scan(driver.Scanner(&u))
		if err != nil {
			return []*manage.ConfigFile{}, fmt.Errorf("failed to scan: %w", err)
		}

		files = append(files, u.Convert())
	}
	if err = stream.Err(); err != nil {
		return []*manage.ConfigFile{}, fmt.Errorf("failed to read from db: %w", err)
	}

	return files, nil
}
