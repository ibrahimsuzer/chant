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

type dotfileRepo interface {
	Add(ctx context.Context, files ...*Dotfile) error
	List(ctx context.Context, page, count int) ([]*Dotfile, error)
}

type dotfileManager struct {
	dotfiles dotfileRepo
}

func NewDotfileManager(dotfileRepo dotfileRepo) *dotfileManager {
	return &dotfileManager{dotfiles: dotfileRepo}
}

func (m *dotfileManager) Add(ctx context.Context, paths ...string) error {

	dotFiles := make([]*Dotfile, 0, len(paths))

	for _, path := range paths {

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

		dotFiles = append(dotFiles, &Dotfile{
			Name:      "",
			Path:      path,
			Extension: filepath.Ext(path),
			MimeType:  mimeType,
			Language:  language,
		})
	}

	err := m.dotfiles.Add(ctx, dotFiles...)
	if err != nil {
		return fmt.Errorf("failed to add config files: %w", err)
	}
	return nil
}

func (m *dotfileManager) List(ctx context.Context) ([]*Dotfile, error) {
	list, err := m.dotfiles.List(ctx, 0, 10)
	if err != nil {
		return []*Dotfile{}, fmt.Errorf("failed to list config files: %w", err)
	}
	return list, nil
}
