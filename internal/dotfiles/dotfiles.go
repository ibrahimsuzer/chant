package dotfiles

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-enry/go-enry/v2"
)

type dotfilePrinter interface {
	Dotfiles(dotfiles ...*Dotfile)
}

type dotfileManager struct {
	dotfiles *dotfileRepo
	printer  dotfilePrinter
}

func NewDotfileManager(dotfileRepo *dotfileRepo, printer dotfilePrinter) *dotfileManager {
	return &dotfileManager{dotfiles: dotfileRepo, printer: printer}
}

func (m *dotfileManager) Add(ctx context.Context, paths ...string) error {

	dotfiles := make([]*Dotfile, 0, len(paths))

	for _, path := range paths {

		absolutePath, err := getAbsolutePath(path)
		if err != nil {
			fmt.Printf("failed to read path: %s", err)

			continue
		}

		// Check details
		content, err := ioutil.ReadFile(absolutePath)
		if err != nil {
			fmt.Printf("failed to read file: %s", err)

			continue
		}

		language := enry.GetLanguage(absolutePath, content)
		mimeType := enry.GetMIMEType(absolutePath, language)

		dotfiles = append(dotfiles, &Dotfile{
			Id:        "",
			Name:      "",
			Path:      absolutePath,
			Extension: filepath.Ext(absolutePath),
			MimeType:  mimeType,
			Language:  language,
		})

	}

	err := m.dotfiles.Add(ctx, dotfiles...)
	if errors.Is(err, ErrUniqueConstraintViolation) {
		fmt.Printf("file already exists")
	} else if err != nil {
		return fmt.Errorf("failed to add config files: %w", err)
	}

	return nil
}

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

func (m *dotfileManager) List(ctx context.Context) error {
	list, err := m.dotfiles.List(ctx, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to list config files: %w", err)
	}

	m.printer.Dotfiles(list...)

	return nil
}
