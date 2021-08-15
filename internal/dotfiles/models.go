package dotfiles

type Dotfile struct {
	Id string

	Name      string
	Path      string
	Extension string
	MimeType  string
	Language  string

	Version *FileVersion
}

type FileVersion struct {
	Id string

	FilePath string
	Version  int
	Content  string
	Hash     string
}
