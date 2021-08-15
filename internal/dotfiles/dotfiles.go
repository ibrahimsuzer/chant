package dotfiles

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-enry/go-enry/v2"
	storage_errors "github.com/ibrahimsuzer/chant/internal/storage/errors"
)

type dotfilePrinter interface {
	Dotfiles(dotfiles ...*Dotfile)
}

type dotfileRepo interface {
	Add(ctx context.Context, file *Dotfile) (*Dotfile, error)
	AddMany(ctx context.Context, files ...*Dotfile) error
	List(ctx context.Context, page, count int) ([]*Dotfile, error)
	Remove(ctx context.Context, ids ...string) error
	Find(ctx context.Context, ids ...string) ([]*Dotfile, error)
}

type fileVersionRepo interface {
	Add(ctx context.Context, fileId string, file *FileVersion) (*FileVersion, error)
}

type dotfileManager struct {
	dotfiles dotfileRepo
	versions fileVersionRepo
	printer  dotfilePrinter
}

func NewDotfileManager(dotfiles dotfileRepo, versions fileVersionRepo, printer dotfilePrinter) *dotfileManager {
	return &dotfileManager{dotfiles: dotfiles, versions: versions, printer: printer}
}


func (m *dotfileManager) Add(ctx context.Context, paths ...string) error {

	for _, path := range paths {

		absolutePath, err := getAbsolutePath(path)
		if err != nil {
			fmt.Printf("failed to read path: %s\n", err)

			continue
		}

		// Check details
		content, err := ioutil.ReadFile(absolutePath)
		if err != nil {
			fmt.Printf("failed to read file: %s\n", err)

			continue
		}

		language := enry.GetLanguage(absolutePath, content)
		mimeType := enry.GetMIMEType(absolutePath, language)

		dotfile, err := m.dotfiles.Add(ctx, &Dotfile{
			Id:        "",
			Name:      "",
			Path:      absolutePath,
			Extension: filepath.Ext(absolutePath),
			MimeType:  mimeType,
			Language:  language,
		})
		if errors.Is(err, storage_errors.ErrUniqueConstraintViolation) {
			fmt.Printf("file already exists: %s \n", absolutePath)
			continue
		} else if err != nil {
			return fmt.Errorf("failed to add dotfile: %w", err)
		}

		_, err = m.versions.Add(ctx, dotfile.Id, &FileVersion{
			Id:      "",
			Content: string(content),
			Hash:    "",
		})
		if err != nil {
			return fmt.Errorf("failed to add file content: %w", err)
		}

	}

	return nil
}

func (m *dotfileManager) List(ctx context.Context) error {
	list, err := m.dotfiles.List(ctx, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to list config files: %w", err)
	}

	m.printer.Dotfiles(list...)

	return nil
}

func (m *dotfileManager) Remove(ctx context.Context, ids ...string) error {

	err := m.dotfiles.Remove(ctx, ids...)
	if err != nil {
		return fmt.Errorf("failed to remove config files: %w", err)
	}

	return nil
}

// Utils

func getAbsolutePath(path string) (string, error) {
	// Check path
	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("path doesn't exist: %s", path)
	} else if err != nil {
		return "", fmt.Errorf("failed to read path: %s", err)
	}

	// Reject directories
	if stat.IsDir() {
		return "", fmt.Errorf("cannot process directories: %s", path)
	}

	// Get absolute path
	evaluatedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("cannot evaluate path: %s", path)
	}

	absolutePath, err := filepath.Abs(evaluatedPath)
	if err != nil {
		return "", fmt.Errorf("cannot find absolute path: %s", path)
	}

	return absolutePath, nil

}
