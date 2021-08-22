package printer

import (
	"github.com/fatih/color"
	"github.com/ibrahimsuzer/chant/internal/dotfiles"
)

type printer struct {
	colorBasic *color.Color
}

func NewPrinter(base *color.Color) *printer {
	return &printer{colorBasic: base}
}

func (p *printer) Dotfiles(dotfiles ...*dotfiles.Dotfile) {
	if len(dotfiles) == 0 {
		p.colorBasic.Printf("Dotfile list is empty\n")

	}
	for _, dotfile := range dotfiles {
		p.colorBasic.Printf("%s %s %s\n", dotfile.Id, dotfile.Name, dotfile.Path)
	}
}
