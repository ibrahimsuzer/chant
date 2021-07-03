package manage


type ConfigFile struct {
	Id           string
	Name         string
	Path         string
	Extension    string
	MimeType     string
	Language     string
}

type ConfigDir struct {
	Id   string
	Name string
	Path string
}
