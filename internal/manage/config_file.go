package manage

type FileFormat string

const (
	YAML FileFormat = "YAML"
	XML  FileFormat = "XML"
	ENV  FileFormat = "ENV"
)

type ConfigFile struct {
	Id          string
	Name        string
	Description string
	Location    string
	Format      FileFormat
}
