package manage

type FileFormat int64

const (
	YAML FileFormat = iota + 1
	XML
	ENV
)

type ConfigFile struct {
	Id          string
	Name        string
	Description string
	Location    string
	Format      FileFormat
}
