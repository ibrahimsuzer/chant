package storage

import (
	"github.com/ibrahimsuzer/chant/internal/manage"
)

type configFile struct {
	Id          string
	Name        string
	Description string
	Location    string
	Format      int64
}

func NewConfigFile(file *manage.ConfigFile, id string) *configFile {
	return &configFile{
		Id:          id,
		Name:        file.Name,
		Description: file.Description,
		Location:    file.Location,
		Format:      int64(file.Format),
	}
}

func (f *configFile) Convert() *manage.ConfigFile {
	return &manage.ConfigFile{
		Id:          f.Id,
		Name:        f.Name,
		Description: f.Description,
		Location:    f.Location,
		Format:      manage.FileFormat(f.Format),
	}
}
