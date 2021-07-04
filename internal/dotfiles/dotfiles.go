package dotfiles

import (
	"context"
	"fmt"
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

func (m *dotfileManager) Add(ctx context.Context, files ...*Dotfile) error {
	err := m.dotfiles.Add(ctx, files...)
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
